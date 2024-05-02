package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"backend/internal/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/lcvvvv/gonmap"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type PortService struct{}

var (
	PortServiceApp = new(PortService)
)

// ExecutePortScan 执行端口扫描任务
func (ps *PortService) ExecutePortScan(c *gin.Context, req request.PortScanRequest, userUUID uuid.UUID) (err error) {
	// 解析 IP 地址 和 端口
	targetList, portList, err := ps.parseRequest(req)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		return err
	}

	// 创建新任务
	task, err := util.StartNewTask(req.Title, strings.Join(targetList, ","), "PortScan", userUUID)
	if err != nil {
		global.Logger.Error("无法创建任务: ", zap.Error(err))
		return errors.New("无法创建任务")
	}

	// 异步执行端口扫描
	go ps.performPortScan(c, req.CheckAlive, task, targetList, portList, req.Threads, req.Timeout)
	return nil
}

func (ps *PortService) parseRequest(portScanRequest request.PortScanRequest) ([]string, []int, error) {
	targetList, err := util.ParseMultipleIPAddresses(portScanRequest.Targets)
	if err != nil {
		return nil, nil, errors.New("IP 地址解析失败")
	}

	if len(targetList) == 0 {
		return nil, nil, errors.New("有效 IP 地址为空")
	}

	portList := util.ParsePort(portScanRequest.Ports)
	if len(portList) == 0 {
		return nil, nil, errors.New("端口解析失败")
	}

	return targetList, portList, nil
}

// performPortScan 执行针对指定目标和端口的端口扫描，使用 ICMP 和自定义的端口扫描逻辑。
func (ps *PortService) performPortScan(c *gin.Context, checkAlive bool, task *model.Task, targets []string, ports []int, threads, timeout int) {
	status := "2" // "2" 表示正在扫描, "3" 表示已取消, "4" 表示错误
	aliveTargets := make(map[string]bool)

	var TargetsMutex sync.Mutex // 互斥锁，保护 Targets
	var statusMutex sync.Mutex  // 互斥锁，保护 status

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, threads) // 用于控制并发数量的信号量

	setStatus := func(s string) {
		statusMutex.Lock()
		defer statusMutex.Unlock()
		status = s
	}

	getStatus := func() string {
		statusMutex.Lock()
		defer statusMutex.Unlock()
		return status
	}

	// 首先检测所有目标是否存活
	if checkAlive {
		for _, target := range targets {
			wg.Add(1)
			semaphore <- struct{}{}
			go func(t string) {
				defer func() {
					<-semaphore
					wg.Done()
				}()
				// 任务取消（主动取消 / 执行失败），停止执行后续方法
				if err := task.Ctx.Err(); err != nil {
					return
				}
				if alive, err := util.IcmpCheckAlive(t, 10); err != nil {
					setStatus("4") // 更新状态为执行失败
					TaskServiceApp.CancelTask(task.UUID, util.GetUUID(c), util.GetAuthorityId(c))
					task.UpdateStatus("4") // 更新任务状态
					global.Logger.Error("检测主机存活失败", zap.String("target", t), zap.Error(err))
					return
				} else if alive {
					TargetsMutex.Lock()
					defer TargetsMutex.Unlock()
					aliveTargets[t] = true
					return
				}
			}(target)
		}
		wg.Wait()

		if err := task.Ctx.Err(); err != nil {
			return
		}
	} else {
		// 如果不检查存活，假定所有目标都是活跃的
		for _, target := range targets {
			aliveTargets[target] = true
		}
	}

	results := make(chan model.PortScanResult, len(ports)*len(aliveTargets)) // 存储扫描结果的通道

	for target, alive := range aliveTargets {
		if alive {
			for _, port := range ports {
				wg.Add(1)
				semaphore <- struct{}{}
				go func(t string, p int) {
					defer func() {
						<-semaphore
						wg.Done()
					}()
					// 任务取消（主动取消 / 执行失败），停止执行后续方法
					if err := task.Ctx.Err(); err != nil {
						return
					}
					//执行端口扫描
					result := ps.PortCheck(t, p, timeout, task.UUID)
					if result != nil {
						results <- *result
					}
				}(target, port)
			}
		}
	}

	// 等待所有扫描任务完成后关闭结果通道
	go func() {
		wg.Wait()
		close(results)
		finalStatus := getStatus()
		task.UpdateStatus(finalStatus)
		if status == "2" { // 扫描正常完成的情况下，收集并处理数据
			ps.processResults(results)
		}
	}()
}

func (ps *PortService) processResults(results chan model.PortScanResult) {
	for result := range results {
		if result.Open {
			err := result.InsertData(global.DB)
			if err != nil {
				global.Logger.Error("插入扫描结果失败: ", zap.Error(err))
			}
		}
	}
}

func (ps *PortService) PortCheck(target string, port int, timeout int, taskUUID uuid.UUID) *model.PortScanResult {
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
	return nil
}

func (ps *PortService) FetchResult(cdb *gorm.DB, result model.PortScanResult, info request.PageInfo, order string, desc bool) ([]model.PortScanResult, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := cdb.Model(&model.PortScanResult{})

	// 条件查询
	if result.TaskUUID != uuid.Nil {
		db = db.Where("task_uuid LIKE ?", "%"+result.TaskUUID.String()+"%")
	}
	if result.IP != "" {
		db = db.Where("ip LIKE ?", "%"+result.IP+"%")
	}
	if result.Service != "" {
		db = db.Where("service LIKE ?", "%"+result.Service+"%")
	}

	// 获取满足条件的条目总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, nil
	}
	// 根据有效列表进行排序处理
	orderStr := "created_at desc" // 默认排序
	if order != "" {
		allowedOrders := map[string]bool{
			"task_uuid":  true,
			"ip":         true,
			"service":    true,
			"created_at": true,
		}
		if _, ok := allowedOrders[order]; !ok {
			return nil, 0, fmt.Errorf("非法的排序字段: %v", order)
		}
		orderStr = order
		if desc {
			orderStr += " desc"
		}
	}

	// 查询数据
	var resultList []model.PortScanResult
	if err := db.Limit(limit).Offset(offset).Order(orderStr).Find(&resultList).Error; err != nil {
		return nil, 0, err
	}

	return resultList, total, nil
}
