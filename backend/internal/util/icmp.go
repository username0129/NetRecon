package util

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"sync"
	"time"
)

func IcmpCheckAlive(target string, timeout int) (bool, error) {
	// 目标地址解析
	dst, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		return false, fmt.Errorf("IP 地址解析错误: %v\n", err)
	}

	// 监听 ICMP
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return false, fmt.Errorf("ICMP 监听失败: %v\n", err)
	}
	defer c.Close()

	// 创建 ICMP Echo 请求消息
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, // 使用进程 ID 和序列号 1
			Data: []byte("HELLO"),
		},
	}

	// 将消息编码为字节
	b, err := msg.Marshal(nil)
	if err != nil {
		return false, fmt.Errorf("消息编码错误: %v\n", err)
	}

	// 发送 ICMP Echo 请求
	if _, err := c.WriteTo(b, dst); err != nil {
		return false, fmt.Errorf("发送 ICMP 请求错误: %v\n", err)
	}

	// 设置超时
	if err = c.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
		return false, fmt.Errorf("设置读取超时错误: %v", err)
	}

	// 接收 ICMP Echo 回应
	reply := make([]byte, 1500)
	n, _, err := c.ReadFrom(reply)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return false, nil // 超时错误，认为目标不存活
		}
		return false, fmt.Errorf("接收回应错误: %v", err)
	}

	// 解析 ICMP 回应
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		return false, fmt.Errorf("解析回应消息错误: %v", err)
	}

	// 检查是否是有效的 ICMP Echo 回复
	if rm.Type == ipv4.ICMPTypeEchoReply {
		fmt.Printf("目标 %v 存活\n", dst.String())
		return true, nil
	}

	return false, nil
}

// checkAllTargetsAlive 检查一组目标的存活状态，限制并发数
func checkAllTargetsAlive(targets []string, threads int, timeout int) {
	semaphore := make(chan struct{}, threads) // 用于控制并发数
	var wg sync.WaitGroup

	for _, target := range targets {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(t string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			alive, err := IcmpCheckAlive(t, timeout)
			if err != nil {
				fmt.Printf("检查 %s 失败: %v\n", t, err)
			} else if alive {
				fmt.Printf("目标 %s 存活\n", t)
			} else {
				fmt.Printf("目标 %s 不存活\n", t)
			}
		}(target)
	}
	wg.Wait()
}
