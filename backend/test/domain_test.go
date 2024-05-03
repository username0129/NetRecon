package test

import (
	"backend/internal/util"
	"fmt"
	"testing"
)

func TestDomainV(t *testing.T) {
	blacklist := []string{
		"baddomain.com",
		"malicious-site.net",
		"example.net",
	}

	testDomains := "example.com,sub.domain.co.uk,-invalid.com, invalid-.com, invalid.com-, 123.com, invalid_site.com, baddomain.com, malicious-site.net, .com, xxx.edu.cn"

	domains, err := util.ParseMultipleDomains(testDomains, blacklist)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("domains: %v\n", domains)
}

func TestDomainAlive(t *testing.T) {
	domain := "example.com"
	if util.IsDomainAlive(domain) {
		fmt.Println(domain, "is alive")
	} else {
		fmt.Println(domain, "is not alive")
	}
}
