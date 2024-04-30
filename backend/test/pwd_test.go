package test

import (
	"backend/internal/util"
	"fmt"
	"testing"
)

func TestPwd(t *testing.T) {
	pwd := util.GetExecPwd()
	fmt.Println(pwd)
}
