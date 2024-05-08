package service

import (
	"backend/internal/config"
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"backend/internal/util"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"time"
)

type CronService struct{}

var (
	CronServiceApp = new(CronService)
)

func (cs *CronService) parseRequest(targets, ports string) ([]string, []int, error) {
	targetList, err := util.ParseMultipleIPAddresses(targets)
	if err != nil {
		return nil, nil, errors.New("IP 地址解析失败")
	}

	if len(targetList) == 0 {
		return nil, nil, errors.New("有效 IP 地址为空")
	}

	portList := util.ParsePort(ports)
	if len(portList) == 0 {
		return nil, nil, errors.New("端口解析失败")
	}

	return targetList, portList, nil
}

func (cs *CronService) AddTask(manager *config.CronManager, req request.CronAddTaskRequest, userUUID uuid.UUID, authorityId string) (err error) {
	runtime := time.Now()
	spec := util.TimeToCronSpec(runtime.Add(30 * time.Second)) // 30 秒之后执行
	if req.TaskType == "PortScan" {
		targetList, err := util.ParseMultipleIPAddresses(req.Targets)
		if err != nil {
			global.Logger.Error("IP 地址解析失败: ", zap.Error(err))
			return errors.New("IP 地址解析失败")
		}

		if len(targetList) == 0 {
			global.Logger.Error("有效 IP 地址为空: ", zap.Error(err))
			return errors.New("有效 IP 地址为空")
		}

		portList := util.ParsePort(req.Ports)
		if len(portList) == 0 {
			global.Logger.Error("端口解析失败: ", zap.Error(err))
			return errors.New("端口解析失败")
		}

		// 创建新任务
		task, err := util.StartNewTask(req.Title, req.Targets, "Cron/Port", req.DictType, userUUID, req.AssetUUID)
		if err != nil {
			global.Logger.Error("无法创建任务: ", zap.Error(err))
			return errors.New("无法创建任务")
		}

		taskID, err := manager.AddTask(spec, createPortScanTaskFunc(req.CheckAlive, task, targetList, portList, req.Threads, req.Timeout, userUUID, authorityId))
		if err != nil {
			global.Logger.Error("添加计划任务失败", zap.Error(err))
			return fmt.Errorf("添加计划任务失败")
		}
		task.CronID = int(taskID)
		task.NextTime = runtime.Add(30 * time.Second).Format("2006-01-02 15:04:05")
		task.CreateOrUpdate()
	} else if req.TaskType == "Subdomain" {

		// 解析域名列表，黑名单校验
		targetList, err := util.ParseMultipleDomains(req.Targets, global.Config.BlackDomain)
		if err != nil {
			global.Logger.Error("域名解析失败: ", zap.String("targets", req.Targets), zap.Error(err))
			return err // 自定义错误
		}

		if len(targetList) == 0 {
			global.Logger.Error("域名解析失败: 有效域名为空")
			return errors.New("有效域名为空")
		}

		// 加载 CDN 列表
		cdnList, err := util.LoadCDNList(util.GetExecPwd() + "/data/cdn.yaml")
		if err != nil {
			global.Logger.Error("加载 CDN 列表失败: ", zap.Error(err))
			return fmt.Errorf("加载 CDN 列表失败")
		}

		// 加载子域名字典
		dict, err := util.LoadSubDomainDict(util.GetExecPwd()+"/data/dict/", req.DictType)
		if err != nil {
			global.Logger.Error("加载子域名字典失败: ", zap.Error(err))
			return fmt.Errorf("加载子域名字典失败")
		}

		// 创建新任务
		task, err := util.StartNewTask(req.Title, req.Targets, "Cron/Domain", req.DictType, userUUID, req.AssetUUID)

		if err != nil {
			global.Logger.Error("无法创建任务: ", zap.String("title", req.Title), zap.Error(err))
			return errors.New("无法创建任务")
		}

		taskID, err := manager.AddTask(spec, createSubdomainTaskFunc(task, targetList, req.Threads, req.Timeout, dict, cdnList, userUUID))
		if err != nil {
			global.Logger.Error("添加计划任务失败", zap.Error(err))
			return fmt.Errorf("添加计划任务失败")
		}
		task.CronID = int(taskID)
		task.NextTime = runtime.Add(30 * time.Second).Format("2006-01-02 15:04:05")
		task.CreateOrUpdate()
	}
	return nil
}

func createPortScanTaskFunc(checkAlive bool, task *model.Task, targets []string, ports []int, threads, timeout int, userUUID uuid.UUID, authorityId string) func() {
	return func() {
		now := time.Now()
		task.LastTime = now.Format("2006-01-02 15:04:05")
		task.NextTime = now.Add(24 * time.Hour).Format("2006-01-02 15:04:05")
		task.CreateOrUpdate()
		PortServiceApp.PerformPortScan(checkAlive, task, targets, ports, threads, timeout, userUUID, authorityId)
	}
}

func createSubdomainTaskFunc(task *model.Task, targets []string, threads, timeout int, dict []string, cdnList map[string][]string, userUUID uuid.UUID) func() {
	return func() {
		now := time.Now()
		task.LastTime = now.Format("2006-01-02 15:04:05")
		task.NextTime = now.Add(24 * time.Hour).Format("2006-01-02 15:04:05")
		task.CreateOrUpdate()
		SubDomainServiceApp.executeBruteSubdomain(task, targets, threads, timeout, dict, cdnList, userUUID)
	}
}
