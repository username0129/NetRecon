package test

import (
	"backend/internal/util"
	"fmt"
	"testing"
)

func TestLoadSubDict(t *testing.T) {
	dict, err := util.LoadSubDomainDict("C:\\Users\\Administrator\\Desktop\\NetRecon\\backend\\dict\\", "large")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dict)
	return
}
