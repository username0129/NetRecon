package model

type AuthorityRoute struct {
	RouteId     int `json:"routeId" gorm:"comment:路由 ID;"`
	AuthorityId int `json:"authorityId" gorm:"comment:角色 ID;"`
}

func (*AuthorityRoute) TableName() string {
	return "sys_authority_route"
}
