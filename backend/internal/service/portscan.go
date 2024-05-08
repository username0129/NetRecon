package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"backend/internal/util"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/lcvvvv/gonmap"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PortService struct{}

var (
	PortServiceApp = new(PortService)
)

// ExecutePortScan 执行端口扫描任务
func (ps *PortService) ExecutePortScan(req request.PortScanRequest, userUUID uuid.UUID, authorityId string, TaskType string) (err error) {
	// 解析 IP 地址 和 端口
	targetList, portList, err := ps.parseRequest(req)
	if err != nil {
		global.Logger.Error("请求解析失败", zap.Error(err))
		return err
	}

	// 创建新任务
	task, err := util.StartNewTask(req.Title, req.Targets, TaskType, req.DictType, userUUID, uuid.Nil)
	if err != nil {
		global.Logger.Error("无法创建任务: ", zap.Error(err))
		return errors.New("无法创建任务")
	}

	// 异步执行端口扫描
	go ps.PerformPortScan(req.CheckAlive, task, targetList, portList, req.Threads, req.Timeout, userUUID, authorityId)
	return nil
}

func (ps *PortService) parseRequest(req request.PortScanRequest) ([]string, []int, error) {
	targetList, err := util.ParseMultipleIPAddresses(req.Targets)
	if err != nil {
		return nil, nil, errors.New("IP 地址解析失败")
	}

	if len(targetList) == 0 {
		return nil, nil, errors.New("有效 IP 地址为空")
	}

	portList := util.ParsePort(req.Ports)
	if len(portList) == 0 {
		return nil, nil, errors.New("端口解析失败")
	}

	return targetList, portList, nil
}

// PerformPortScan 执行针对指定目标和端口的端口扫描，使用 ICMP 和自定义的端口扫描逻辑。
func (ps *PortService) PerformPortScan(checkAlive bool, task *model.Task, targets []string, ports []int, threads, timeout int, userUUID uuid.UUID, authorityId string) {
	status := "1" // "1" 表示正在扫描, "2" 表示扫描完成, "3" 表示已取消, "4" 表示错误
	aliveTargets := make(map[string]bool)

	var statusMutex sync.Mutex // 互斥锁，保护 status
	var TargetMutex sync.Mutex // 互斥锁，保护 aliveTargets

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
			go func(t string) {
				defer wg.Done()
				defer func() { <-semaphore }()
				semaphore <- struct{}{}
				// 任务状态出现变动，如取消 / 执行失败
				if getStatus() != "1" {
					return
				}
				if task.Status == "3" {
					setStatus("3") // 更新状态为取消
					return
				}
				if alive, err := util.IcmpCheckAlive(t, 10); err != nil {
					setStatus("4") // 更新状态为执行失败
					TaskServiceApp.CancelTask(task.UUID, userUUID, authorityId)
					task.UpdateStatus("4") // 更新任务状态
					global.Logger.Error("检测主机存活失败", zap.String("target", t), zap.Error(err))
				} else if alive {
					TargetMutex.Lock()
					defer TargetMutex.Unlock()
					aliveTargets[t] = true
				}
			}(target)
		}
		wg.Wait()
	} else {
		// 如果不检查存活，假定所有目标都是活跃的
		for _, target := range targets {
			aliveTargets[target] = true
		}
	}

	if getStatus() != "1" {
		return
	}

	results := make(chan model.PortScanResult, len(ports)*len(aliveTargets)) // 存储扫描结果的通道

	for target, alive := range aliveTargets {
		if alive {
			for _, port := range ports {
				wg.Add(1)
				go func(t string, p int) {
					defer wg.Done()
					defer func() { <-semaphore }()
					semaphore <- struct{}{}
					// 任务状态出现变动，如取消 / 执行失败
					if getStatus() != "1" {
						return
					}
					if task.Status == "3" {
						setStatus("3") // 更新状态为取消
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
		if finalStatus == "1" { // 扫描正常完成的情况下，收集并处理数据
			task.UpdateStatus("2")
			ps.processResults(results, userUUID, task)
		}
	}()
}

func (ps *PortService) processResults(results chan model.PortScanResult, userUUID uuid.UUID, task *model.Task) {
	count := 0
	for result := range results {
		if result.Open {
			err := result.InsertData(global.DB)
			if err != nil {
				global.Logger.Error("插入扫描结果失败: ", zap.Error(err))
			} else {
				count++
			}
		}
	}

	userMail, err := UserServiceApp.GetUserMailByUUID(userUUID)
	if err != nil {
		global.Logger.Error("获取用户邮箱失败: ", zap.Error(err))
	} else {
		timeCompleted := time.Now().Format("2006-01-02 15:04:05")
		body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: 'Arial', sans-serif; line-height: 1.6; }
    h1 { color: #333; }
	p { margin: 10px 0; }
    .footer { color: grey; font-size: 0.9em; }
    hr { border: 0; height: 1px; background-color: #ddd; }
  </style>
</head>
<body>
  <h1>任务执行完成通知</h1>
  <p><strong>任务标题：</strong>%s</p>
  <p><strong>目标：</strong>%s</p>
  <p><strong>完成时间：</strong>%s</p>
  <p><strong>获得有效数据：</strong>%d 条</p>
  <hr>
  <p class="footer">此邮件为系统自动发送，请勿直接回复。</p>
</body>
</html>
`, task.Title, task.Targets, timeCompleted, count)
		subject := fmt.Sprintf("端口扫描任务完成通知 - UUID %s", task.UUID)
		mail := global.Config.Mail
		err := util.SendMail(mail.SmtpServer, mail.SmtpPort, mail.SmtpFrom, mail.SmtpPassword, userMail, subject, body)
		if err != nil {
			global.Logger.Error("发送邮箱失败: ", zap.Error(err))
		}
	}
}

func (ps *PortService) PortCheck(target string, port int, timeout int, taskUUID uuid.UUID) *model.PortScanResult {
	result := model.PortScanResult{
		UUID:     uuid.Must(uuid.NewV4()),
		TaskUUID: taskUUID,
		IP:       target,
		Port:     port,
	}
	scanner := gonmap.New()
	scanner.SetTimeout(time.Duration(timeout) * time.Second)
	status, response := scanner.Scan(target, port)
	switch status {
	case gonmap.Closed:
		result.Open = false
	case gonmap.Unknown:
		result.Open = true
		result.Service = "filter" // filter 未知状态
	default:
		result.Open = true
	}
	if response != nil {
		if strings.TrimSpace(response.FingerPrint.Service) != "" {
			result.Service = response.FingerPrint.Service
		} else {
			result.Service = "unknown"
		}
	}
	return &result
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
	if result.Port != 0 {
		db = db.Where("port LIKE ?", "%"+strconv.Itoa(result.Port)+"%")
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
	orderStr := "ip desc" // 默认排序
	if order != "" {
		allowedOrders := map[string]bool{
			"ip":           true,
			"port":         true,
			"creator_uuid": true,
			"service":      true,
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

// DeleteResult 删除端口扫描结果
func (ps *PortService) DeleteResult(uuid uuid.UUID) error {
	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 直接尝试删除记录
	result := tx.Where("uuid = ?", uuid).Delete(&model.PortScanResult{})
	if result.Error != nil {
		tx.Rollback() // 如果删除操作出错，回滚事务
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback() // 如果没有删除任何记录，回滚事务
		return errors.New("没有找到记录")
	}

	// 提交事务
	return tx.Commit().Error
}

// DeleteResults  删除端口扫描结果
func (ps *PortService) DeleteResults(uuids []uuid.UUID) error {
	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 直接尝试删除记录
	result := tx.Where("uuid in (?)", uuids).Delete(&model.PortScanResult{})
	if result.Error != nil {
		tx.Rollback() // 如果删除操作出错，回滚事务
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback() // 如果没有删除任何记录，回滚事务
		return errors.New("没有找到记录")
	}

	// 提交事务
	return tx.Commit().Error
}

// FetchAllResult 获取全部数据
func (ps *PortService) FetchAllResult(db *gorm.DB, taskUUID uuid.UUID) ([]model.PortScanResult, error) {
	var result []model.PortScanResult
	if err := db.Model(&model.PortScanResult{}).Where("task_uuid LIKE ?", "%"+taskUUID.String()+"%").Find(&result).Error; err != nil {
		global.Logger.Error("查询数据失败: ", zap.Error(err))
		// 查询数据失败
		return nil, errors.New("查询数据失败")
	}
	return result, nil
}
