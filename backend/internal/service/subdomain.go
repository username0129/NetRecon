package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"backend/internal/util"
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gofrs/uuid/v5"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SubDomainService struct {
}

var (
	SubDomainServiceApp = new(SubDomainService)
)

// BruteSubdomains 爆破子域名
func (ss *SubDomainService) BruteSubdomains(subdomainRequest request.SubDomainRequest, userUUID uuid.UUID) (err error) {
	// 创建新任务
	task, err := util.StartNewTask(subdomainRequest.Title, subdomainRequest.Targets, "BruteSubdomain", userUUID)
	if err != nil {
		global.Logger.Error("无法创建任务: ", zap.String("title", subdomainRequest.Title), zap.Error(err))
		return errors.New("无法创建任务")
	}

	targetList, err := util.ParseMultipleDomains(subdomainRequest.Targets, global.Config.BlackDomain)
	if err != nil {
		global.Logger.Error("域名解析失败: ", zap.String("targets", subdomainRequest.Targets), zap.Error(err))
		task.UpdateStatus("4")
		return err
	}

	if len(targetList) == 0 {
		global.Logger.Error("域名解析失败: 有效域名为空")
		task.UpdateStatus("4")
		return errors.New("有效域名为空")
	}

	dict, err := util.LoadSubDomainDict(util.GetExecPwd()+"/data/dict/", subdomainRequest.DictType)
	if err != nil {
		global.Logger.Error("加载子域名字典失败: ", zap.Error(err))
		task.UpdateStatus("4")
		return err
	}

	cdnList, err := util.LoadCDNList(util.GetExecPwd() + "/data/cdn.yaml")
	if err != nil {
		global.Logger.Error("加载 CDN 文件失败: ", zap.Error(err))
		task.UpdateStatus("4")
		return err
	}

	threads, err := strconv.Atoi(subdomainRequest.Threads)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		task.UpdateStatus("4")
		return errors.New("线程数必须是整数")
	}

	timeout, err := strconv.Atoi(subdomainRequest.Timeout)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		task.UpdateStatus("4")
		return errors.New("超时时间必须是整数")
	}

	go ss.executeBruteSubdomain(task, targetList, threads, timeout, dict, cdnList)

	return nil
}

func (ss *SubDomainService) executeBruteSubdomain(task *model.Task, targets []string, threads int, timeout int, dict []string, cdnList map[string][]string) {
	status := "2"              // "2" 表示正在扫描, "3" 表示已取消, "4" 表示错误
	var statusMutex sync.Mutex // 互斥锁，保护 status
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, threads) // 用于控制并发数量的信号量

	setStatus := func(s string) {
		statusMutex.Lock()
		defer statusMutex.Unlock()
		status = s
	}

	getStatus := func() string {
		statusMutex.Lock()
		defer statusMutex.Unlock()
		return status
	}

	results := make(chan model.SubDomainResult, len(targets)*len(dict)) // 存储扫描结果的通道
	for _, target := range targets {
		for _, sub := range dict {
			wg.Add(1)
			semaphore <- struct{}{}
			go func(t, sub string) {
				defer func() {
					<-semaphore
					wg.Done()
				}()

				subdomain := sub + "." + t
				if getStatus() != "2" {
					return
				}

				// 解析域名 CNAME 和 IP
				result, err := ss.Resolution(task.Ctx, subdomain, timeout, task.UUID, cdnList)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						setStatus("3") // 更新状态为已取消
					} else {
						if !(errors.Is(err, errors.New("signal: killed")) && !errors.Is(err, errors.New("no such host"))) {
							setStatus("4") // 更新状态为出错
							global.Logger.Error("检测域名失败", zap.String("target", subdomain), zap.Error(err))
						}
					}
				}
				if err == nil && result != nil {
					results <- *result
					fmt.Printf("域名: %v, CNAME: %v, IPS: %v\n", (*result).SubDomain, (*result).Cname, (*result).Ips)
				}
			}(target, sub)
		}
	}

	// 等待所有扫描任务完成后关闭结果通道
	go func() {
		wg.Wait()
		close(results)
		finalStatus := getStatus()
		task.UpdateStatus(finalStatus)
		if status == "2" { // 扫描正常完成的情况下，收集并处理数据
			ss.processResults(results)
		}
	}()

}

// Resolution 解析域名
func (ss *SubDomainService) Resolution(ctx context.Context, domain string, timeout int, taskUUID uuid.UUID, cdnList map[string][]string) (subDomainResult *model.SubDomainResult, err error) {
	cnames, err := ss.LookupCNAME(ctx, domain, timeout)
	if err != nil {
		return nil, err
	}

	ips, err := ss.LookupHost(ctx, domain, timeout)
	if err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, nil
	}

	url := "http://" + domain
	title, code, err := ss.FetchTitle(url)
	if err != nil {
		global.Logger.Error("获取网站标题失败", zap.String("url", domain), zap.Error(err))
		if strings.Contains(err.Error(), "connection refused") {
			title = " NONE / 站点拒绝连接"
		} else if code != 0 {
			title = fmt.Sprintf("%v / 状态码: %v", code, code)
		}
	} else {
		title = "200 / " + strings.TrimSpace(title)
	}

	subDomainResult = &model.SubDomainResult{
		TaskUUID:  taskUUID,
		SubDomain: domain,
		Title:     title,
		Cname:     strings.Join(cnames, ","),
		Ips:       strings.Join(ips, ","),
	}

	for _, cdns := range cdnList {
		for _, cdn := range cdns {
			for _, cname := range cnames {
				if strings.Contains(cname, cdn) { // 识别到cdn
					subDomainResult.Notes = fmt.Sprintf("在 CNAME 中识别到CDN字段%v", cdn)
					return subDomainResult, nil
				} else if strings.Contains(cname, "cdn") {
					subDomainResult.Notes = fmt.Sprintf("在 CNAME %v中检测到 cdn 关键字", cname)
					return subDomainResult, nil
				}
			}
		}
	}
	return subDomainResult, nil
}

// LookupCNAME 获取 DNS CNAME 记录 （判断 CDN）
func (ss *SubDomainService) LookupCNAME(ctx context.Context, domain string, timeout int) (CNAMES []string, err error) {
	c := dns.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeCNAME)

	r, _, err := c.ExchangeContext(ctx, &m, "114.114.114.114:53")
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, nil // 超时
		} else {
			return nil, err
		}
	}

	if len(r.Answer) == 0 {
		return nil, nil // 没有 CNAME 记录
	}

	for _, ans := range r.Answer {
		if record, isType := ans.(*dns.CNAME); isType {
			CNAMES = append(CNAMES, record.Target)
		}
	}
	return CNAMES, nil
}

// LookupHost 查找主机 IP 地址
func (ss *SubDomainService) LookupHost(ctx context.Context, domain string, timeout int) (hosts []string, err error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(timeout) * time.Second,
			}
			return d.DialContext(ctx, "tcp", "114.114.114.114:53")
		},
	}

	ips, err := r.LookupHost(ctx, domain)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			if strings.Contains(dnsErr.Err, "server misbehaving") {
				//global.Logger.Warn("DNS 服务器 misbehaving error ignored", zap.String("domain", domain), zap.Error(err))
				return nil, nil // 忽略错误，返回空结果
			}
			if dnsErr.IsNotFound || dnsErr.IsTimeout {
				return nil, nil // 无法解析到 IP，返回空结果，忽略错误
			}
		}
		return nil, err // 其他错误，正常处理
	} else {
		return ips, nil
	}
}

func (ss *SubDomainService) processResults(results chan model.SubDomainResult) {
	for result := range results {
		err := result.InsertData(global.DB)
		if err != nil {
			global.Logger.Error("插入扫描结果失败: ", zap.Error(err))
		}
	}
}

// FetchTitle 发送 HTTP GET 请求到指定的 URL 并解析 HTML 文档以提取网页标题。
func (ss *SubDomainService) FetchTitle(url string) (string, int, error) {
	// 创建 HTTP 客户端对象，允许重定向
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 返回最后一个响应，不自动重定向
		},
	}

	// 设置请求头部，模拟浏览器访问
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("获取 URL 失败: %w", err)
	}

	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther {
		// 如果是重定向响应，使用 Location 头部中的 URL 重新发起请求
		if location, err := resp.Location(); err == nil && location.String() != "" {
			return ss.FetchTitle(location.String()) // 递归调用处理重定向
		}
	} else if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, fmt.Errorf("无效的响应码: %d", resp.StatusCode)
	}

	// 加载 HTML 文档
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("加载 HTML 响应失败: %w", err)
	}

	// 提取 <title> 标签的内容
	title := doc.Find("title").Text()
	return title, 0, nil
}
