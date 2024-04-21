package model

type AuthorityRoute struct {
	RouteId     uint `json:"routeId" gorm:"comment:路由 ID;"`
	AuthorityId uint `json:"authorityId" gorm:"comment:角色 ID;"`
}

func (*AuthorityRoute) TableName() string {
	return "sys_authority_route"
}
