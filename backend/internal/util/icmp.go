package util

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
)

func IcmpCheckAlive(target string, timeout int) (bool, error) {

	// 目标地址解析
	dst, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		return false, fmt.Errorf("IP 地址解析错误: %v\n", err)
	}

	// 监听本地 ICMP
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return false, fmt.Errorf("ICMP 监听失败: %v\n", err)
	}
	defer conn.Close()

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
	if _, err := conn.WriteTo(b, dst); err != nil {
		return false, fmt.Errorf("发送 ICMP 请求错误: %v\n", err)
	}

	// 设置超时
	if err = conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
		return false, fmt.Errorf("设置读取超时错误: %v", err)
	}

	// 接收 ICMP Echo 回应
	reply := make([]byte, 1500)
	n, _, err := conn.ReadFrom(reply)
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
		return true, nil
	}

	return false, nil
}
