package util

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"time"
)

func IcmpCheckAlive(target string) (bool, error) {

	// 监听 ICMP 数据包
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return false, fmt.Errorf("ICMP 监听失败: %v", err)
	}
	defer conn.Close()

	// 创建 ICMP Echo 请求消息
	id := rand.Intn(0xffff)
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   id,
			Seq:  1,
			Data: []byte("HELLO"),
		},
	}

	// 将消息编组为二进制形式
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return false, fmt.Errorf("编码 ICMP 消息时出错: %v", err)
	}

	// 解析 IP 地址并发送 ICMP 回显请求
	dst, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		return false, fmt.Errorf("解析 IP 地址时出错: %v", err)
	}

	_, err = conn.WriteTo(msgBytes, dst)
	if err != nil {
		return false, fmt.Errorf("发送 ICMP 消息时出错: %v", err)
	}

	// 设置读取超时时间并等待响应
	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return false, fmt.Errorf("设置读取超时时出错: %v", err)
	}

	// 读取 ICMP Echo 响应消息
	reply := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return false, nil
	}

	// 解析 ICMP 回应
	resp, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return false, fmt.Errorf("解析 ICMP 消息时出错: %v", err)
	}

	// 判断响应消息的类型
	switch resp.Type {
	case ipv4.ICMPTypeEchoReply:
		// 检查回显应答的源地址是否与目标地址匹配
		if echo, ok := resp.Body.(*icmp.Echo); ok && echo.ID == id && peer.String() == target {
			return true, nil
		}
	default:
		return false, nil
	}

	return false, nil
}
