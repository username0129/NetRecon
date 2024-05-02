package util

import (
	"strconv"
	"strings"
)

// ParsePort 解析表示端口号的字符串，支持单个端口、逗号分隔的端口列表和连字符表示的端口范围。
func ParsePort(ports string) (portList []int) {
	if ports == "" {
		return nil
	}
	slices := strings.Split(ports, ",")
	for _, port := range slices {
		port = strings.TrimSpace(port)
		if port == "" {
			continue
		}

		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) != 2 {
				continue // 如果范围格式不正确，则跳过
			}

			// 尝试转换范围的两端
			startPort, err1 := strconv.Atoi(ranges[0])
			endPort, err2 := strconv.Atoi(ranges[1])
			if err1 != nil || err2 != nil || startPort > endPort || startPort < 1 || endPort > 65535 {
				continue // 如果转换失败或端口号不在有效范围内，则跳过
			}

			// 将范围内的所有端口加入到列表中
			for i := startPort; i <= endPort; i++ {
				portList = append(portList, i)
			}
		} else { // 单个端口
			if portNum, err := strconv.Atoi(port); err == nil && portNum >= 1 && portNum <= 65535 {
				portList = append(portList, portNum)
			}
		}
	}

	return RemoveDuplicates[int](portList) // 移除重复的端口号并返回结果
}
