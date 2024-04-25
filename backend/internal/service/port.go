package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/util"
	"context"
	"errors"
	"github.com/Ullaakut/nmap/v3"
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
func (ps *PortService) ExecutePortScan(portScanRequest model.PortScanRequest) (err error) {
	// 创建新任务
	task, err := util.StartNewTask("PortScan")
	if err != nil {
		global.Logger.Error("无法创建任务: ", zap.Error(err))
		return errors.New("无法创建任务")
	}

	// 解析 IP 地址 和 端口
	targetList, portList, err := ps.parseRequest(portScanRequest)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		task.UpdateStatus("执行失败")
		return err
	}

	threads, err := strconv.Atoi(portScanRequest.Threads)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		task.UpdateStatus("执行失败")
		return errors.New("线程数必须是整数")
	}

	timeout, err := strconv.Atoi(portScanRequest.Timeout)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		task.UpdateStatus("执行失败")
		return errors.New("超时时间必须是整数")
	}

	// 异步执行端口扫描
	go ps.performPortScan(task.Ctx, task, targetList, portList, threads, timeout)
	return nil
}

func (ps *PortService) parseRequest(portScanRequest model.PortScanRequest) ([]string, []string, error) {
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

func (ps *PortService) performPortScan(ctx context.Context, task *model.Task, targets, ports []string, threads, timeout int) {
	semaphore := make(chan struct{}, threads)
	results := make(chan model.PortScanResult, len(ports)*len(targets))
	var wg sync.WaitGroup

	for _, target := range targets {
		for _, port := range ports {
			wg.Add(1)
			go func(t, p string) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				result, err := ps.PortScan(ctx, t, p, timeout)
				if err == nil && result != nil {
					results <- *result
				}
			}(target, port)
		}
	}

	// 等待所有扫描任务完成后关闭结果通道
	go func() {
		wg.Wait()
		close(results)
		ps.processResults(task, results)
	}()
}

func (ps *PortService) processResults(task *model.Task, results chan model.PortScanResult) {
	openFound := false
	for result := range results {
		if result.Open {
			openFound = true
			err := result.InsertData(global.DB)
			if err != nil {
				global.Logger.Error("插入扫描结果失败: ", zap.Error(err))
			}
		}
	}
	finalStatus := "任务完成"
	if !openFound {
		finalStatus = "未找到开放端口"
	}
	task.UpdateStatus(finalStatus)
}

// PortScan 扫描端口
func (ps *PortService) PortScan(ctx context.Context, target string, port string, timeout int) (*model.PortScanResult, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	result := model.PortScanResult{}

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
