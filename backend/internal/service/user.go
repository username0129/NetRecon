package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"backend/internal/util"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserService struct{}

var (
	UserServiceApp = new(UserService)
)

// AddUserInfo  添加用户
func (us *UserService) AddUserInfo(db *gorm.DB, user model.User) (err error) {
	var count int64
	if db.Model(model.User{}).Where("username = ?", user.Username).Count(&count).Error != nil {
		return fmt.Errorf("查询用户数据失败")
	}
	if count > 0 {
		return fmt.Errorf("用户名'%s' 已存在", user.Username) // 用户名已存在
	} else {
		user.Password = util.BcryptHash(user.Password)
		return user.InsertData(db)
	}
}

// UpdateUserInfo  更新用户信息
func (us *UserService) UpdateUserInfo(db *gorm.DB, user model.User) (err error) {
	// 更新用户信息
	result := db.Model(&model.User{}).Where("uuid = ?", user.UUID).Updates(&user)
	if result.Error != nil {
		return fmt.Errorf("更新用户信息失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("没有找到要更新的用户")
	}
	return nil
}

// UpdatePasswordInfo  更新用户密码
func (us *UserService) UpdatePasswordInfo(db *gorm.DB, req request.UpdatePasswordRequest, userUUID uuid.UUID) (err error) {
	var user model.User
	if err = global.DB.Model(&model.User{}).Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return errors.New("查询用户失败")
	}

	if ok := util.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("旧密码错误")
	}
	user.Password = util.BcryptHash(req.NewPassword)
	// 更新用户信息
	result := db.Model(&model.User{}).Where("uuid = ?", user.UUID).Updates(&user)
	if result.Error != nil {
		return fmt.Errorf("更新用户密码失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("没有找到要更新的用户")
	}
	return nil
}

// ResetPassword  重置用户密码为随机 10 个字符
func (us *UserService) ResetPassword(db *gorm.DB, userUUID uuid.UUID) (err error) {
	userInfo, err := us.FetchUserByUUID(userUUID)
	if err != nil {
		return fmt.Errorf("未找到用户信息")
	}

	newPassword, _ := util.GeneratePassword(10)

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
  <h1>NetRecon 密码重置通知</h1>
  <p><strong>账户：</strong>%s</p>
  <p><strong>新密码：</strong>%s</p>
  <p><strong>重置时间：</strong>%s</p>
  <hr>
  <p class="footer">此邮件为系统自动发送，请勿直接回复。</p>
</body>
</html>
`, userInfo.Username, newPassword, timeCompleted)

	subject := fmt.Sprintf("端口扫描任务完成通知 - 账户 %v", userInfo.Username)
	mail := global.Config.Mail
	err = util.SendMail(mail.SmtpServer, mail.SmtpPort, mail.SmtpFrom, mail.SmtpPassword, []string{userInfo.Mail}, subject, body)
	if err != nil {
		global.Logger.Error("发送重置密码邮箱失败: ", zap.Error(err))
		return fmt.Errorf("发送重置密码邮箱失败，终止重置密码")
	}

	// 更新用户信息
	result := db.Model(&model.User{}).Where("uuid = ?", userUUID).UpdateColumn("password", util.BcryptHash(newPassword))
	if result.Error != nil {
		return fmt.Errorf("重置用户密码失败")
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("没有找到要重置密码的用户")
	}
	return nil
}

// DeleteUserInfo  删除用户
func (us *UserService) DeleteUserInfo(db *gorm.DB, userUUID uuid.UUID) (err error) {

	var taskList []model.Task
	// 先获取用户创建的任务
	if err = db.Model(model.Task{}).Where("creator_uuid = ?", userUUID).Find(&taskList).Error; err != nil {
		return fmt.Errorf("获取用户创建的任务失败")
	}

	// 检查是否有进行中的任务
	for _, task := range taskList {
		if task.Status == "1" {
			return fmt.Errorf("请等待该用户的所有扫描任务执行完成后再进行删除。")
		}
	}

	// 删除所有任务
	for _, task := range taskList {
		if err = TaskServiceApp.DeleteTask(db, task.UUID, userUUID, "1"); err != nil {
			return fmt.Errorf("删除用户创建的任务失败")
		}
	}

	// 删除指定 UUID 的用户
	result := db.Model(model.User{}).Where("uuid = ?", userUUID).Delete(&model.User{})

	// 检查错误
	if result.Error != nil {
		return fmt.Errorf("删除用户数据失败")
	}

	// 检查是否有行被删除
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户名不存在")
	}

	return nil
}

// FetchUserByUUID  根据 UUID 获取用户详细信息
func (us *UserService) FetchUserByUUID(userUUID uuid.UUID) (user model.User, err error) {
	if err := global.DB.Model(&model.User{}).Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return model.User{}, err
	} else {
		return user, nil
	}
}

// FetchUsers 获取用户信息
func (us *UserService) FetchUsers(cdb *gorm.DB, result model.User, info request.PageInfo, order string, desc bool) ([]model.User, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := cdb.Model(&model.User{})

	// 条件查询
	if result.UUID != uuid.Nil {
		db = db.Where("uuid LIKE ?", "%"+result.UUID.String()+"%")
	}
	if result.Username != "" {
		db = db.Where("username LIKE ?", "%"+result.Username+"%")
	}
	if result.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+result.Nickname+"%")
	}
	if result.AuthorityId != "" {
		db = db.Where("authority_id LIKE ?", "%"+result.AuthorityId+"%")
	}
	if result.Enable != "" {
		db = db.Where("enable LIKE ?", "%"+result.Enable+"%")
	}
	if result.Mail != "" {
		db = db.Where("mail LIKE ?", "%"+result.Mail+"%")
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
			"uuid":         true,
			"username":     true,
			"nickname":     true,
			"authority_id": true,
			"enable":       true,
			"mail":         true,
			"created_at":   true,
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
	var resultList []model.User
	if err := db.Limit(limit).Offset(offset).Order(orderStr).Find(&resultList).Error; err != nil {
		return nil, 0, err
	}

	return resultList, total, nil
}

// GetAdministratorMail 获取管理员邮箱
func (us *UserService) GetAdministratorMail() (mails []string, err error) {
	if err := global.DB.Model(&model.User{}).Select("mail").Where("authority_id = 1").Find(&mails).Error; err != nil {
		return nil, err
	} else {
		return mails, nil
	}
}

// GetUserMailByUUID 根据 UUID 获取用户邮箱
func (us *UserService) GetUserMailByUUID(userUUID uuid.UUID) (mails []string, err error) {
	if err := global.DB.Model(&model.User{}).Select("mail").Where("uuid = ?", userUUID).Find(&mails).Error; err != nil {
		return nil, err
	} else {
		return mails, nil
	}
}
