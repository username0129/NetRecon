package test

import (
	"context"
	"fmt"
	"github.com/lcvvvv/gonmap"
	"sync"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	// 设置目标主机和端口范围
	target := "121.37.217.131"
	portRanges := []int{21, 22, 80, 81, 135, 139, 443, 445, 1433, 1521, 3306, 5432, 6379, 7001, 8000, 8080, 8089, 9000, 9200, 11211, 27017, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 98, 99, 443, 800, 801, 808, 880, 888, 889, 1000, 1010, 1080, 1081, 1082, 1099, 1118, 1888, 2008, 2020, 2100, 2375, 2379, 3000, 3008, 3128, 3505, 5555, 6080, 6648, 6868, 7000, 7001, 7002, 7003, 7004, 7005, 7007, 7008, 7070, 7071, 7074, 7078, 7080, 7088, 7200, 7680, 7687, 7688, 7777, 7890, 8000, 8001, 8002, 8003, 8004, 8006, 8008, 8009, 8010, 8011, 8012, 8016, 8018, 8020, 8028, 8030, 8038, 8042, 8044, 8046, 8048, 8053, 8060, 8069, 8070, 8080, 8081, 8082, 8083, 8084, 8085, 8086, 8087, 8088, 8089, 8090, 8091, 8092, 8093, 8094, 8095, 8096, 8097, 8098, 8099, 8100, 8101, 8108, 8118, 8161, 8172, 8180, 8181, 8200, 8222, 8244, 8258, 8280, 8288, 8300, 8360, 8443, 8448, 8484, 8800, 8834, 8838, 8848, 8858, 8868, 8879, 8880, 8881, 8888, 8899, 8983, 8989, 9000, 9001, 9002, 9008, 9010, 9043, 9060, 9080, 9081, 9082, 9083, 9084, 9085, 9086, 9087, 9088, 9089, 9090, 9091, 9092, 9093, 9094, 9095, 9096, 9097, 9098, 9099, 9100, 9200, 9443, 9448, 9800, 9981, 9986, 9988, 9998, 9999, 10000, 10001, 10002, 10004, 10008, 10010, 10250, 12018, 12443, 14000, 16080, 18000, 18001, 18002, 18004, 18008, 18080, 18082, 18088, 18090, 18098, 19001, 20000, 20720, 21000, 21501, 21502, 28018, 20880} // 示例端口范围

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // 用于控制并发数量的信号量

	// 创建一个可以取消的context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 在2秒后自动取消扫描作为示例
	time.AfterFunc(1*time.Second, cancel)

	for _, port := range portRanges {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(t string, p int) {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			if err := ctx.Err(); err != nil {
				fmt.Println("新任务添加被阻止，上下文已取消")
				return
			}
			runScan(t, p)
		}(target, port)
	}

	wg.Wait()
	fmt.Println("所有扫描任务已完成或被取消")
}

func runScan(target string, port int) {
	scanner := gonmap.New()
	status, response := scanner.Scan(target, port)
	switch status {
	case gonmap.Closed:
		fmt.Printf("%v:%v %v", target, port, "closed")
	// filter 未知状态
	case gonmap.Unknown:
		fmt.Printf("%v:%v %v", target, port, "unknown")
	default:
		fmt.Printf("%v:%v %v ", target, port, "open")
	}
	if response != nil {
		if response.FingerPrint.Service != "" {
			fmt.Printf("Service: %v\n", response.FingerPrint.Service)
		} else {
			fmt.Printf("Service: unknown\n")
		}
	}
}
