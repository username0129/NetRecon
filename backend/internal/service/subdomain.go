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
	"gorm.io/gorm"
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
func (ss *SubDomainService) BruteSubdomains(req request.SubDomainRequest, userUUID uuid.UUID, TaskType string) (err error) {
	// 解析域名列表，黑名单校验
	targetList, err := util.ParseMultipleDomains(req.Targets, global.Config.BlackDomain)
	if err != nil {
		global.Logger.Error("域名解析失败: ", zap.String("targets", req.Targets), zap.Error(err))
		return err // 自定义错误
	}

	if len(targetList) == 0 {
		global.Logger.Error("域名解析失败: 有效域名为空")
		return errors.New("有效域名为空")
	}

	// 加载 CDN 列表
	cdnList, err := util.LoadCDNList(util.GetExecPwd() + "/data/cdn.yaml")
	if err != nil {
		global.Logger.Error("加载 CDN 列表失败: ", zap.Error(err))
		return fmt.Errorf("加载 CDN 列表失败")
	}

	// 加载子域名字典
	dict, err := util.LoadSubDomainDict(util.GetExecPwd()+"/data/dict/", req.DictType)
	if err != nil {
		global.Logger.Error("加载子域名字典失败: ", zap.Error(err))
		return fmt.Errorf("加载子域名字典失败")
	}

	// 创建新任务
	task, err := util.StartNewTask(req.Title, req.Targets, TaskType, req.DictType, userUUID, uuid.Nil)

	if err != nil {
		global.Logger.Error("无法创建任务: ", zap.String("title", req.Title), zap.Error(err))
		return errors.New("无法创建任务")
	}

	go ss.executeBruteSubdomain(task, targetList, req.Threads, req.Timeout, dict, cdnList, userUUID)

	return nil
}

func (ss *SubDomainService) executeBruteSubdomain(task *model.Task, targets []string, threads, timeout int, dict []string, cdnList map[string][]string, userUUID uuid.UUID) {
	status := "1"              // "2" 表示正在扫描, "3" 表示已取消, "4" 表示错误
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
			go func(t, sub string) {
				semaphore <- struct{}{}
				defer wg.Done()
				defer func() { <-semaphore }()
				// 任务状态出现变动，如取消 / 执行失败
				if getStatus() != "1" {
					return
				}
				if task.Status == "3" {
					setStatus("3") // 更新状态为取消
					return
				}
				subdomain := sub + "." + t
				// 解析域名 CNAME 和 IP
				result, err := ss.Resolution(task.Ctx, subdomain, timeout, task.UUID, cdnList)
				if err != nil {
					if errors.Is(err, context.Canceled) || strings.Contains(err.Error(), "operation was canceled") {
						setStatus("3") // 更新状态为已取消
					} else {
						//if !(errors.Is(err, errors.New("signal: killed")) && !errors.Is(err, errors.New("no such host"))) {
						setStatus("4") // 更新状态为出错
						global.Logger.Error("检测域名失败", zap.String("target", subdomain), zap.Error(err))
						//}
					}
				}
				if err == nil && result != nil {
					results <- *result
				}
			}(target, sub)
		}
	}

	// 等待所有扫描任务完成后关闭结果通道
	go func() {
		wg.Wait()
		close(results)
		finalStatus := getStatus()
		if finalStatus == "1" { // 扫描正常完成的情况下，收集并处理数据
			task.UpdateStatus("2") // 更新任务状态为 2 -> 扫描完成
			ss.processResults(results, userUUID, task)
		}
	}()

}

// Resolution 解析域名
func (ss *SubDomainService) Resolution(ctx context.Context, domain string, timeout int, taskUUID uuid.UUID, cdnList map[string][]string) (subDomainResult *model.SubDomainResult, err error) {
	// 解析域名 CNAME 列表
	cnames, err := ss.LookupCNAME(ctx, domain, timeout)
	if err != nil {
		return nil, err
	}

	// 解析域名指向的 IP 地址
	ips, err := ss.LookupHost(ctx, domain, timeout)
	if err != nil {
		return nil, err
	}

	// 当 IP 地址指向为空时说明域名不存活，跳出后续网页标题获取
	if len(ips) == 0 {
		return nil, nil
	}

	// 先使用 http 协议进行探测
	url := "http://" + domain
	title, code, err := ss.FetchTitle(url)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			title = "站点拒绝连接"
		} else {
			if strings.Contains(err.Error(), "timeout") {
				title = "站点连接超时"
			} else {
				title = "站点连接失败"
			}
		}
	}

	subDomainResult = &model.SubDomainResult{
		UUID:      uuid.Must(uuid.NewV4()),
		TaskUUID:  taskUUID,
		SubDomain: domain,
		Title:     strings.TrimSpace(title),
		Code:      code,
		Cname:     strings.Join(cnames, ","),
		Ips:       ips[0],
	}

	for _, cdns := range cdnList {
		for _, cdn := range cdns {
			for _, cname := range cnames {
				if strings.Contains(cname, cdn) { // 识别到cdn
					subDomainResult.Notes = fmt.Sprintf("在 CNAME 中识别到 CDN 字段: %v", cdn)
					return subDomainResult, nil
				} else if strings.Contains(cname, "cdn") {
					subDomainResult.Notes = fmt.Sprintf("在 CNAME 中检测到 cdn 关键字")
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

	ips, err := r.LookupIPAddr(ctx, domain)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			if strings.Contains(dnsErr.Err, "server misbehaving") {
				return nil, nil // 忽略错误，返回空结果
			}
			if dnsErr.IsNotFound || dnsErr.IsTimeout {
				return nil, nil // 无法解析到 IP，返回空结果，忽略错误
			}
		}
		return nil, err
	}

	// 筛选出 IPv4 地址
	var ipv4s []string
	for _, ipAddr := range ips {
		if ipv4 := ipAddr.IP.To4(); ipv4 != nil {
			ipv4s = append(ipv4s, ipv4.String())
		}
	}

	return ipv4s, nil
}

func (ss *SubDomainService) processResults(results chan model.SubDomainResult, userUUID uuid.UUID, task *model.Task) {
	count := 0
	for result := range results {
		err := result.InsertData(global.DB)
		if err != nil {
			global.Logger.Error("插入扫描结果失败: ", zap.Error(err))
		} else {
			count++
		}
	}

	userMail, err := UserServiceApp.GetUserMailByUUID(userUUID)
	if err != nil {
		global.Logger.Error("获取用户邮箱失败: ", zap.Error(err))
	} else {
		timeCompleted := time.Now().Format("2006-01-02 15:04:05")
		body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: 'Arial', sans-serif; line-height: 1.6; }
    h1 { color: #333; }
	p { margin: 10px 0; }
    .footer { color: grey; font-size: 0.9em; }
    hr { border: 0; height: 1px; background-color: #ddd; }
  </style>
</head>
<body>
  <h1>任务执行完成通知</h1>
  <p><strong>任务标题：</strong>%s</p>
  <p><strong>目标：</strong>%s</p>
  <p><strong>完成时间：</strong>%s</p>
  <p><strong>获得有效数据：</strong>%d 条</p>
  <hr>
  <p class="footer">此邮件为系统自动发送，请勿直接回复。</p>
</body>
</html>
`, task.Title, task.Targets, timeCompleted, count)
		subject := fmt.Sprintf("子域名扫描任务完成通知 - UUID %s", task.UUID)
		mail := global.Config.Mail
		err := util.SendMail(mail.SmtpServer, mail.SmtpPort, mail.SmtpFrom, mail.SmtpPassword, userMail, subject, body)
		if err != nil {
			global.Logger.Error("发送邮箱失败: ", zap.Error(err))
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
		// 如果存在重定向响应，如 https，使用 Location 头部中的 URL 重新发起请求
		if location, err := resp.Location(); err == nil && location.String() != "" {
			return ss.FetchTitle(location.String()) // 递归调用处理重定向
		}
	}

	// 加载 HTML 响应文档
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("加载 HTML 响应失败: %w", err)
	}

	// 提取 <title> 标签的内容
	title := doc.Find("title").Text()

	// 当标题为无效的 UTF8 编码时，尝试进行清洗无效字符
	if !util.ValidateUTF8(title) {
		title = util.CleanString(title)
	}
	return title, resp.StatusCode, nil
}

func (ss *SubDomainService) FetchResult(cdb *gorm.DB, result model.SubDomainResult, info request.PageInfo, order string, desc bool) ([]model.SubDomainResult, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := cdb.Model(&model.SubDomainResult{})
	// 条件查询
	if result.TaskUUID != uuid.Nil {
		db = db.Where("task_uuid LIKE ?", "%"+result.TaskUUID.String()+"%")
	}
	if result.SubDomain != "" {
		db = db.Where("sub_domain LIKE ?", "%"+result.SubDomain+"%")
	}
	if result.Title != "" {
		db = db.Where("title LIKE ?", "%"+result.Title+"%")
	}
	if result.Ips != "" {
		db = db.Where("ips LIKE ?", "%"+result.Ips+"%")
	}
	if result.Cname != "" {
		db = db.Where("cname LIKE ?", "%"+result.Cname+"%")
	}
	if result.Code != 0 {
		db = db.Where("code LIKE ?", "%"+strconv.Itoa(result.Code)+"%")
	}

	// 获取满足条件的条目总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, nil
	}
	// 根据有效列表进行排序处理
	orderStr := "sub_domain desc" // 默认排序
	if order != "" {
		allowedOrders := map[string]bool{
			"task_uuid":    true,
			"sub_domain":   true,
			"title":        true,
			"ips":          true,
			"code":         true,
			"creator_uuid": true,
			"created_at":   true,
		}
		if _, ok := allowedOrders[order]; !ok {
			return nil, 0, fmt.Errorf("非法的排序字段: %v", order)
		}
		orderStr = order
		if desc {
			orderStr += " desc"
		}
	}

	// 查询数据
	var resultList []model.SubDomainResult
	if err := db.Limit(limit).Offset(offset).Order(orderStr).Find(&resultList).Error; err != nil {
		return nil, 0, err
	}

	return resultList, total, nil
}
