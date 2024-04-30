package util

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"net"
	"strconv"
	"strings"
)

// parseIPAddress 解析单个 IP 地址或范围
func parseIPAddress(input string) ([]string, error) {
	var ips []string

	// CIDR 格式，如 192.168.0.1/24
	if strings.Contains(input, "/") {
		_, network, err := net.ParseCIDR(input)
		if err != nil {
			return nil, err
		}
		for ip := network.IP.Mask(network.Mask); network.Contains(ip); incIP(ip) {
			ips = append(ips, ip.String())
		}
		// 移除网络和广播地址
		if len(ips) > 2 {
			ips = ips[1 : len(ips)-1]
		}
		return ips, nil
	}

	// Range 格式，如 192.168.0.1-255
	if strings.Contains(input, "-") {
		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format")
		}
		startIP, endOctetStr := parts[0], parts[1]
		base, err := net.ResolveIPAddr("ip", startIP)
		if err != nil {
			return nil, err
		}
		start := base.IP
		endOctet, err := strconv.Atoi(endOctetStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end of range: %s", parts[1])
		}
		for i := start[len(start)-1]; i <= byte(endOctet); i++ {
			newIP := net.ParseIP(start.String())
			newIP[len(newIP)-1] = i
			ips = append(ips, newIP.String())
		}
		return ips, nil
	}

	// 单个 IP
	ip := net.ParseIP(input)
	if ip != nil {
		ips = append(ips, ip.String())
		return ips, nil
	}

	return nil, fmt.Errorf("invalid IP format")
}

// incIP 从最后一位开始递增
// net.IP 为 byte 类型，范围为 0-255，当加 255 + 1 = 0。如果最后以为 0 则说明当前 /24 网段已经扫描结束将会往前一段遍历
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// ParseMultipleIPAddresses 处理包含多种格式的 IP 地址的字符串
func ParseMultipleIPAddresses(input string) ([]string, error) {
	var allIPs []string
	parts := strings.Split(input, ",")
	for _, part := range parts {
		ips, err := parseIPAddress(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		allIPs = append(allIPs, ips...)
	}
	// 去重
	return RemoveDuplicates[string](allIPs), nil
}

func Ip2region(ip string) (string, error) {
	var dbPath = "./data/ip2region/ip2region.xdb"
	searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		return "", err
	}
	defer searcher.Close()
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		return "", err
	}
	return region, nil
}
