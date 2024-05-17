package controller

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/common"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/service"
	"backend/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	JWTRequired bool
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	uuid := util.GetUUID(c)
	if global.DB == nil {
		common.ResponseError(c, http.StatusInternalServerError, "数据库未初始化", nil)
		return
	}

	if user, err := service.UserServiceApp.FetchUserByUUID(uuid); err != nil {
		global.Logger.Error("获取用户信息失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, "获取用户信息失败", nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "获取用户信息成功", user)
	}
}

func (uc *UserController) PostResetPassword(c *gin.Context) {
	var req request.UUIDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostResetPassword 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	if err := service.UserServiceApp.ResetPassword(global.DB, req.UUID); err != nil {
		global.Logger.Error("重置用户密码失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("重置用户密码失败: %v", err.Error()), nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "重置用户密码成功", nil)
	}
}

func (uc *UserController) PostUpdateUserInfo(c *gin.Context) {
	var req request.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostUpdateUserInfo 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	user := model.User{
		UUID:        req.UUID,
		Username:    req.Username,
		Nickname:    req.Nickname,
		Mail:        req.Mail,
		Avatar:      req.Avatar,
		AuthorityId: req.AuthorityId,
		Enable:      req.Enable,
	}

	err := service.UserServiceApp.UpdateUserInfo(global.DB, user)
	if err != nil {
		common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("更新用户失败: %v", err), nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "更新用户成功", nil)
	return
}

func (uc *UserController) PostFetchUsers(c *gin.Context) {
	var req request.FetchUsersRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFetchUsers 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	result, total, err := service.UserServiceApp.FetchUsers(global.DB, req.User, req.PageInfo, req.OrderKey, req.Desc)

	if err != nil {
		global.Logger.Error("查询数据失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, "查询数据失败", nil)
		return
	}

	if total == 0 {
		common.ResponseOk(c, http.StatusNotFound, "未查询到有效数据", nil)
		return
	} else {
		common.ResponseOk(c, http.StatusOK, "查询数据成功", response.PageResult{
			Data:     result,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		})
		return
	}
}

func (uc *UserController) PostAddUserInfo(c *gin.Context) {
	var req request.AddUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostAddUserInfo 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	user := model.User{
		Username:    req.Username,
		Password:    req.Password,
		Nickname:    req.Nickname,
		Mail:        req.Mail,
		Avatar:      req.Avatar,
		AuthorityId: req.AuthorityId,
		Enable:      req.Enable,
	}

	err := service.UserServiceApp.AddUserInfo(global.DB, user)
	if err != nil {
		global.Logger.Error("添加用户失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("添加用户失败: %v", err), nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "添加用户成功", nil)
	return
}

// PostDeleteUserInfo 删除用户数据
func (uc *UserController) PostDeleteUserInfo(c *gin.Context) {
	var req request.UUIDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostFetchUsers 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	err := service.UserServiceApp.DeleteUserInfo(global.DB, req.UUID)
	if err != nil {
		global.Logger.Error("删除用户失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("删除用户失败: %v", err), nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "删除用户成功", nil)
	return
}

// PostUpdatePassword 更新用户密码
func (uc *UserController) PostUpdatePassword(c *gin.Context) {
	var req request.UpdatePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("PostUpdatePassword 参数解析错误: ", zap.Error(err))
		common.ResponseOk(c, http.StatusBadRequest, "参数解析错误", nil)
		return
	}

	err := service.UserServiceApp.UpdatePasswordInfo(global.DB, req, util.GetUUID(c))
	if err != nil {
		global.Logger.Error("更新用户密码失败: ", zap.Error(err))
		common.ResponseOk(c, http.StatusInternalServerError, fmt.Sprintf("更新用户密码失败: %v", err), nil)
		return
	}
	common.ResponseOk(c, http.StatusOK, "删除用户成功", nil)
	return
}
