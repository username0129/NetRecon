package test

import (
	"backend/internal/util"
	"fmt"
	"testing"
)

func TestPassword(t *testing.T) {
	fmt.Println(util.BcryptHash("123456"))
	fmt.Println(util.BcryptCheck("123456", "$2a$10$KHrd/nV2/p4P5YP1B02bOOUhBIHm/ypNfE9QfaY2rp/E0tkvsFX3S"))
}
