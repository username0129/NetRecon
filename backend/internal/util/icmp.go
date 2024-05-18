package util

import (
	"backend/internal/global"
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	"time"
)

// ICMP Echo请求和回应的类型和代码
const (
	icmpEchoRequestType = 8 // ICMP 请求
	icmpEchoReplyType   = 0
	icmpEchoRequestCode = 0 // ICMP 应答
	icmpEchoReplyCode   = 0
	timeoutDuration     = 2 * time.Second
)

// ICMP 消息结构
type ICMP struct {
	Type        uint8  // 消息类型，这里只使用 8 和 0
	Code        uint8  // 通常为 0
	Checksum    uint16 // 校验和
	Identifier  uint16 // 发送方标识符
	SequenceNum uint16 // 请求序列号
}

// 计算ICMP数据包的校验和
func checksum(data []byte) uint16 {
	var sum uint32

	// 2 字节相加
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}
	// 最后一个字节右补零
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	// 溢出的值加到低位上
	for sum>>16 != 0 {
		sum = (sum >> 16) + (sum & 0xFFFF)
	}
	return uint16(^sum)
}

func IcmpCheckAlive(target string) bool {
	// 创建原始socket，IPv4及ICMP协议
	conn, err := net.Dial("ip4:icmp", target)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("和 %v 建立 ICMP 连接失败：", target), zap.Error(err))
		return false
	}
	defer conn.Close()

	icmp := ICMP{
		Type:        icmpEchoRequestType, // Echo请求
		Code:        icmpEchoRequestCode,
		Checksum:    0, // 初始校验和为0
		Identifier:  12345,
		SequenceNum: 1,
	}

	data := make([]byte, 8) // 创建足够的空间存储ICMP头部
	binary.BigEndian.PutUint16(data[0:], uint16(icmp.Type)<<8+uint16(icmp.Code))
	binary.BigEndian.PutUint16(data[2:], uint16(0))
	binary.BigEndian.PutUint16(data[4:], icmp.Identifier)
	binary.BigEndian.PutUint16(data[6:], icmp.SequenceNum)

	// 计算校验和
	icmp.Checksum = checksum(data)
	binary.BigEndian.PutUint16(data[2:4], icmp.Checksum)

	// 发送ICMP请求
	if _, err := conn.Write(data); err != nil {
		fmt.Println("发送 ICMP 请求失败：", err)
		os.Exit(1)
	}

	// 接收ICMP响应
	reply := make([]byte, 20+len(data)) // 包括IP头部
	done := make(chan bool, 1)
	go func() {
		_, err := conn.Read(reply)
		if err != nil {
			done <- false
			return
		}
		icmpType := reply[20] // 解析ICMP类型
		if icmpType == icmpEchoReplyType {
			done <- true
		} else {
			done <- false
		}
	}()

	select {
	case success := <-done:
		if success {
			return true
		} else {
			return false
		}
	case <-time.After(timeoutDuration):
		return false
	}
}
