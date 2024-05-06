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
	"github.com/robfig/cron/v3"
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

func (cs *CronService) AddTask(manager *config.CronManager, req request.CronAddTaskRequest, userUUID uuid.UUID, authorityId string) (taskID cron.EntryID, err error) {
	runtime := time.Now()
	spec := util.TimeToCronSpec(runtime.Add(30 * time.Second)) // 30 秒之后执行
	if req.TaskType == "PortScan" {
		targetList, err := util.ParseMultipleIPAddresses(req.Targets)
		if err != nil {
			global.Logger.Error("IP 地址解析失败: ", zap.Error(err))
			return taskID, errors.New("IP 地址解析失败")
		}

		if len(targetList) == 0 {
			global.Logger.Error("有效 IP 地址为空: ", zap.Error(err))
			return taskID, errors.New("有效 IP 地址为空")
		}

		portList := util.ParsePort(req.Ports)
		if len(portList) == 0 {
			global.Logger.Error("端口解析失败: ", zap.Error(err))
			return taskID, errors.New("端口解析失败")
		}

		// 创建新任务
		task, err := util.StartNewTask(req.Title, req.Targets, "Cron/Port", req.DictType, userUUID, req.AssetUUID)
		if err != nil {
			global.Logger.Error("无法创建任务: ", zap.Error(err))
			return taskID, errors.New("无法创建任务")
		}

		taskID, err = manager.AddTask(spec, createPortScanTaskFunc(req.CheckAlive, task, targetList, portList, req.Threads, req.Timeout, userUUID, authorityId))
		if err != nil {
			global.Logger.Error("添加计划任务失败", zap.Error(err))
			return taskID, fmt.Errorf("添加计划任务失败")
		}
		task.CronID = int(taskID)
		task.CreateOrUpdate()
	} else if req.TaskType == "BruteSubdomain" {
		subdomainRequest := request.SubDomainRequest{
			Title:    req.Title,
			Targets:  req.Targets,
			Timeout:  req.Timeout,
			Threads:  req.Threads,
			DictType: req.DictType,
		}

		taskID, err = manager.AddTask(spec, createSubdomainTaskFunc(subdomainRequest, userUUID))
		if err != nil {
			global.Logger.Error("添加计划任务失败", zap.Error(err))
			return taskID, fmt.Errorf("添加计划任务失败")
		}
	}
	return taskID, nil
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

func createSubdomainTaskFunc(req request.SubDomainRequest, userUUID uuid.UUID) func() {
	return func() {
		err := SubDomainServiceApp.BruteSubdomains(req, userUUID, "CronTask/BrutSubdomain")
		if err != nil {
			global.Logger.Error("计划任务执行失败", zap.Error(err))
		}
	}
}
