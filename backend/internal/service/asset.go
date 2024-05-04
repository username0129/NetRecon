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
	"strings"
)

type AssetService struct{}

var (
	AssetServiceApp = new(AssetService)
)

func (as *AssetService) parseRequest(req request.AddAssetRequest) ([]string, []string, error) {
	ipList, err := util.ParseMultipleIPAddresses(req.IPs)
	if err != nil {
		return nil, nil, errors.New("IP 地址解析失败")
	}

	if len(ipList) == 0 {
		return nil, nil, errors.New("有效 IP 地址为空")
	}

	domainList, err := util.ParseMultipleDomains(req.Domains, global.Config.BlackDomain)
	if err != nil {
		global.Logger.Error("域名解析失败: ", zap.String("targets", req.Domains), zap.Error(err))
		return nil, nil, err // 自定义错误
	}

	return ipList, domainList, nil
}

// DomainExists 检查域名是否已经存在
func (as *AssetService) DomainExists(db *gorm.DB, domain string) (bool, error) {
	var count int64
	if err := db.Model(&model.Asset{}).Where("domains LIKE ?", "%"+domain+"%").Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IPExists 检查IP是否已经存在
func (as *AssetService) IPExists(db *gorm.DB, ip string) (bool, error) {
	var count int64
	if err := db.Model(&model.Asset{}).Where("ips LIKE ?", "%"+ip+"%").Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (as *AssetService) CheckExists(db *gorm.DB, domainList, ipList []string) error {
	// 判断域名是否存在重叠
	for _, domain := range domainList {
		exists, err := as.DomainExists(db, domain)
		if err != nil {
			global.Logger.Error("判断域名重叠发生错误", zap.String("domain", domain), zap.Error(err))
			return fmt.Errorf("检查域名 '%s' 时发生错误: %v", domain, err)
		}
		if exists {
			global.Logger.Warn("域名已存在", zap.String("domain", domain))
			return fmt.Errorf("域名 '%s' 已存在", domain)
		}
	}

	// 判断 IP 是否存在重叠
	for _, ip := range ipList {
		exists, err := as.IPExists(db, ip)
		if err != nil {
			global.Logger.Error("判断 IP 重叠发生错误", zap.String("ip", ip), zap.Error(err))
			return fmt.Errorf("检查 IP '%s' 时发生错误: %v", ip, err)
		}
		if exists {
			global.Logger.Warn("IP 已存在", zap.String("ip", ip))
			return fmt.Errorf("IP '%s' 已存在", ip)
		}
	}
	return nil
}

func (as *AssetService) AddAsset(db *gorm.DB, req request.AddAssetRequest, userUUID uuid.UUID) (err error) {
	// 解析列表
	ipList, domainList, err := as.parseRequest(req)
	if err != nil {
		return err
	}
	// 判断资产是否存在重叠
	if err = as.CheckExists(db, domainList, ipList); err != nil {
		return err
	}
	// 没有重叠的情况下，插入到数据库中
	asset := model.Asset{
		UUID:        uuid.Must(uuid.NewV4()),
		CreatorUUID: userUUID,
		Domains:     strings.Join(domainList, ","),
		IPs:         strings.Join(ipList, ","),
	}

	if err = asset.InsertData(db); err != nil {
		return fmt.Errorf("资产提交失败")
	}
	return nil
}
