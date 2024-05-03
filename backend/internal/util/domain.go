package util

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func IsDomainAlive(domain string) bool {
	_, err := net.LookupHost(domain)
	return err == nil
}

// ParseMultipleDomains 解析域名
func ParseMultipleDomains(input string, blackDomain []string) (allDomains []string, err error) {
	parts := strings.Split(input, ",")
	for _, part := range parts {
		domain := strings.TrimSpace(part)
		if len(domain) == 0 {
			continue
		}
		validDomain, err := isValidDomain(domain, blackDomain)
		if err != nil {
			return nil, err
		}
		if validDomain {
			allDomains = append(allDomains, domain)
		}
	}
	// 去重
	return RemoveDuplicates[string](allDomains), nil
}

// 判断域名是否合法
func isValidDomain(domain string, blackDomain []string) (bool, error) {
	// 正则表达式用于匹配合法的域名，忽略大小写，判断一 a-z,0-9 开头以 . 结尾的部分，最后以 a-z,0-9 结尾的部分
	var domainRegex = regexp.MustCompile(`^(?i)([a-z0-9]([-a-z0-9]*[a-z0-9])?\.)*[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	// 检查域名长度是否符合规定
	if len(domain) > 253 {
		return false, fmt.Errorf("域名过长")
	}

	// 检查域名是否在黑名单中
	for _, b := range blackDomain {
		if strings.HasSuffix(domain, b) {
			return false, fmt.Errorf("域名 %v 在黑名单中", domain)
		}
	}

	// 检查域名是否符合格式
	if !domainRegex.MatchString(domain) {
		return false, fmt.Errorf("无效的域名: %s", domain)
	}

	return true, nil
}
