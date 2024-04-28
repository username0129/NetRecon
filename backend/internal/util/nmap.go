package util

import (
	"backend/internal/model"
	"context"
	"github.com/Ullaakut/nmap/v3"
	"github.com/gofrs/uuid/v5"
	"time"
)

// CheckHostAlive 通过 ping 来探测目标是否存活
func CheckHostAlive(ctx context.Context, target string, timeout int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// 创建扫描器
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(target), // 目标
		nmap.WithPingScan(),
	)

	if err != nil {
		return false, err
	}

	scanResult, _, err := scanner.Run()
	if err != nil {
		return false, err
	}
	if len(scanResult.Hosts) > 0 {
		if scanResult.Hosts[0].Status.State == "up" {
			return true, nil
		}
	}
	return false, nil // 扫描成功但没有开放
}

func PortScan(ctx context.Context, target string, port string, timeout int, taskUUID uuid.UUID) (*model.PortScanResult, error) {
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
