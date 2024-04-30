package test

import (
	"backend/internal/util"
	"fmt"
	"testing"
)

func TestIp2Region(t *testing.T) {
	region, err := util.Ip2region("121.37.217.131")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Region:", region)
}
