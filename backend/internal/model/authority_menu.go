package model

type AuthorityMenu struct {
	MenuId      string `json:"menuId" gorm:"comment:菜单 ID;"`
	AuthorityId string `json:"authorityId" gorm:"comment:角色 ID;"`
}

func (*AuthorityMenu) TableName() string {
	return "sys_authority_menu"
}
