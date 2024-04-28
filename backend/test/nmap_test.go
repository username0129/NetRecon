package test

import (
	"backend/internal/util"
	"context"
	"testing"
)

func TestNmapPing(t *testing.T) {
	_, _ = util.CheckHostAlive(context.Background(), "192.168.80.1", 20)
}
