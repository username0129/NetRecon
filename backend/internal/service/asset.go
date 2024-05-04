package service

import (
	"backend/internal/global"
	"backend/internal/model"
	"backend/internal/model/request"
	"backend/internal/util"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type AssetService struct{}

var (
	AssetServiceApp = new(AssetService)
)

func (as *AssetService) AddAsset(db *gorm.DB, req request.AddAssetRequest, userUUID uuid.UUID) (err error) {
	if len(req.IPs) != 0 {
		// 解析 IP 地址
		ipList, err := util.ParseMultipleIPAddresses(req.IPs)
		if err != nil {
			return errors.New("IP 地址解析失败")
		}

		if len(ipList) == 0 {
			return errors.New("有效 IP 地址为空")
		}
	}

	if len(req.Domains) != 0 {
		// 解析域名列表
		domainList, err := util.ParseMultipleDomains(req.Domains, global.Config.BlackDomain)
		if err != nil {
			return err // 自定义错误
		}

		if len(domainList) == 0 {
			return errors.New("有效域名为空")
		}
	}

	//插入到数据库中
	asset := model.Asset{
		UUID:        uuid.Must(uuid.NewV4()),
		CreatorUUID: userUUID,
		Title:       req.Title,
		Domains:     req.Domains,
		IPs:         req.IPs,
	}
	if err = asset.InsertData(db); err != nil {
		return fmt.Errorf("资产提交失败")
	}
	return nil
}

// UpdateAsset  更新资产信息
func (as *AssetService) UpdateAsset(db *gorm.DB, req request.UpdateAssetRequest) (err error) {

	if len(req.IPs) != 0 {
		// 解析 IP 地址
		ipList, err := util.ParseMultipleIPAddresses(req.IPs)
		if err != nil {
			return errors.New("IP 地址解析失败")
		}

		if len(ipList) == 0 {
			return errors.New("有效 IP 地址为空")
		}
	}

	if len(req.Domains) != 0 {
		// 解析域名列表
		domainList, err := util.ParseMultipleDomains(req.Domains, global.Config.BlackDomain)
		if err != nil {
			return err // 自定义错误
		}

		if len(domainList) == 0 {
			return errors.New("有效域名为空")
		}
	}

	asset := model.Asset{
		UUID:    req.UUID,
		Title:   req.Title,
		Domains: req.Domains,
		IPs:     req.IPs,
	}

	// 更新用户信息，确保零值更新
	result := db.Model(&model.Asset{}).Where("uuid = ?", req.UUID).Select("UUID", "Title", "Domains", "IPs").Updates(asset)
	if result.Error != nil {
		return fmt.Errorf("更新资产信息失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("没有找到要更新的资产")
	}
	return nil
}

// DeleteAsset  删除资产
func (as *AssetService) DeleteAsset(db *gorm.DB, userUUID uuid.UUID) (err error) {

	// 删除指定 UUID 的用户
	result := db.Model(&model.Asset{}).Where("uuid = ?", userUUID).Delete(&model.Asset{})

	// 检查错误
	if result.Error != nil {
		return fmt.Errorf("删除资产数据失败")
	}

	// 检查是否有行被删除
	if result.RowsAffected == 0 {
		return fmt.Errorf("资产不存在")
	}

	return nil
}

func (as *AssetService) FetchAssets(cdb *gorm.DB, req model.Asset, info request.PageInfo, order string, desc bool, userUUID uuid.UUID, authorityId string) ([]model.Asset, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := cdb.Model(&model.Asset{})

	// 管理员用户可查看全部扫描任务
	if authorityId != "1" {
		db = db.Where("creator_uuid LIKE ?", "%"+userUUID.String()+"%")
	}

	// 条件查询
	if req.UUID != uuid.Nil {
		db = db.Where("uuid LIKE ?", "%"+req.UUID.String()+"%")
	}
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Domains != "" {
		db = db.Where("domains LIKE ?", "%"+req.Domains+"%")
	}
	if req.IPs != "" {
		db = db.Where("ips LIKE ?", "%"+req.IPs+"%")
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
			"title":        true,
			"domains":      true,
			"ips":          true,
			"created_at":   true,
			"creator_uuid": true,
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
	var resultList []model.Asset
	if err := db.Preload("Creator").Limit(limit).Offset(offset).Order(orderStr).Find(&resultList).Error; err != nil {
		return nil, 0, err
	}

	return resultList, total, nil
}
