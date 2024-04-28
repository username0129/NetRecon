package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"backend/internal/util"
	"context"
	"errors"
	"github.com/Ullaakut/nmap/v3"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type PortService struct{}

var (
	PortServiceApp = new(PortService)
)

// ExecutePortScan 执行端口扫描任务
func (ps *PortService) ExecutePortScan(portScanRequest request.PortScanRequest, userUUID uuid.UUID) (err error) {
	// 创建新任务
	task, err := util.StartNewTask(portScanRequest.Title, "PortScan", userUUID)
	if err != nil {
		global.Logger.Error("无法创建任务: ", zap.Error(err))
		return errors.New("无法创建任务")
	}

	// 解析 IP 地址 和 端口
	targetList, portList, err := ps.parseRequest(portScanRequest)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		task.UpdateStatus("4")
		return err
	}

	threads, err := strconv.Atoi(portScanRequest.Threads)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		task.UpdateStatus("4")
		return errors.New("线程数必须是整数")
	}

	timeout, err := strconv.Atoi(portScanRequest.Timeout)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		task.UpdateStatus("4")
		return errors.New("超时时间必须是整数")
	}

	// 异步执行端口扫描
	go ps.performPortScan(task.Ctx, task, targetList, portList, threads, timeout, task.UUID)
	return nil
}

func (ps *PortService) parseRequest(portScanRequest request.PortScanRequest) ([]string, []string, error) {
	targetList, err := util.ParseMultipleIPAddresses(portScanRequest.Targets)
	if err != nil || len(targetList) == 0 {
		return nil, nil, errors.New("IP 地址解析失败")
	}

	portList := util.ParsePort(portScanRequest.Ports)
	if len(portList) == 0 {
		return nil, nil, errors.New("端口解析失败")
	}

	return targetList, portList, nil
}

// performPortScan 执行针对指定目标和端口的端口扫描，使用 ICMP 和自定义的端口扫描逻辑。
func (ps *PortService) performPortScan(ctx context.Context, task *model.Task, targets, ports []string, threads, timeout int, taskUUID uuid.UUID) {
	var status string = "2" // "2" 表示正在扫描, "3" 表示已取消, "4" 表示错误
	aliveTargets := make(map[string]bool)
	var TargetsMutex sync.Mutex // 互斥锁，保护 Targets
	var statusMutex sync.Mutex  // 互斥锁，保护 status

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, threads) // 用于控制并发数量的信号量

	// 首先检测所有目标是否存活
	for _, target := range targets {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(t string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if alive, _ := util.IcmpCheckAlive(t, 5); alive {
				TargetsMutex.Lock()
				aliveTargets[t] = true
				TargetsMutex.Unlock()
			}
		}(target)
	}
	wg.Wait()

	results := make(chan model.PortScanResult, len(ports)*len(aliveTargets)) // 存储扫描结果的通道

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

	for target, alive := range aliveTargets {
		if alive {
			for _, port := range ports {
				wg.Add(1)
				semaphore <- struct{}{}
				go func(t, p string) {
					defer wg.Done()
					defer func() { <-semaphore }()

					if getStatus() != "2" {
						return
					}

					//执行端口扫描
					result, err := ps.PortScan(ctx, t, p, timeout, taskUUID)
					if err != nil {
						if errors.Is(err, context.Canceled) {
							setStatus("3") // 更新状态为已取消
						} else {
							setStatus("4") // 更新状态为出错
							global.Logger.Error("端口扫描错误", zap.String("target", t), zap.String("port", p), zap.Error(err))
						}
					}
					if err == nil && result != nil {
						results <- *result
					}
				}(target, port)
			}
		}
	}

	// 等待所有扫描任务完成后关闭结果通道
	go func() {
		wg.Wait()
		close(results)
		finalStatus := getStatus()
		task.UpdateStatus(finalStatus)
		if status == "2" { // 扫描正常完成的情况下，收集并处理数据
			ps.processResults(results)
		}
	}()
}

func (ps *PortService) processResults(results chan model.PortScanResult) {
	for result := range results {
		if result.Open {
			err := result.InsertData(global.DB)
			if err != nil {
				global.Logger.Error("插入扫描结果失败: ", zap.Error(err))
			}
		}
	}
}

// PortScan 扫描端口
func (ps *PortService) PortScan(ctx context.Context, target string, port string, timeout int, taskUUID uuid.UUID) (*model.PortScanResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	result := model.PortScanResult{
		TaskUUID: taskUUID,
	}

	// 创建扫描器
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(target),     // 目标
		nmap.WithPorts(port),         // 端口
		nmap.WithServiceInfo(),       // 识别服务
		nmap.WithSkipHostDiscovery(), // 跳过存活探测
	)

	if err != nil {
		return nil, err
	}

	scanResult, _, err := scanner.Run()
	if err != nil {
		return nil, err
	}

	if len(scanResult.Hosts) > 0 {
		host := scanResult.Hosts[0]
		if len(host.Ports) > 0 {
			port := host.Ports[0]
			result.IP = host.Addresses[0].String()
			result.Port = port.ID
			result.Protocol = port.Protocol
			result.Open = port.State.State == "open"
			result.Service = port.Service.Name
			return &result, nil // 成功返回结果
		}
	}

	return &result, nil // 扫描成功但没有找到开放的端口
}
