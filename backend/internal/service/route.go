package service

import (
	"backend/internal/global"
	"backend/internal/model"
)

type RouteService struct{}

var (
	RouterServiceApp = new(RouteService)
)

func (rs *RouteService) GetRouteTree(authorityId string) (baseRoute []model.Route, err error) {
	var routes []model.Route

	var authorityRoutes []model.AuthorityRoute
	if err = global.DB.Model(&model.AuthorityRoute{}).Where("authority_id = ?", authorityId).Find(&authorityRoutes).Error; err != nil {
		return nil, err
	}

	var routeIds []int
	for _, item := range authorityRoutes {
		routeIds = append(routeIds, item.RouteId)
	}

	if err = global.DB.Model(&model.Route{}).Where("id in (?)", routeIds).Find(&routes).Error; err != nil {
		return nil, err
	}

	routeMap := make(map[uint][]model.Route)

	for _, route := range routes {
		routeMap[route.ParentId] = append(routeMap[route.ParentId], route)
	}

	baseRoute = routeMap[0]
	for i := 0; i < len(baseRoute); i++ {
		rs.GetChildrenList(&baseRoute[i], routeMap)
	}

	return baseRoute, nil
}

func (rs *RouteService) GetChildrenList(route *model.Route, routeMap map[uint][]model.Route) {
	if route != nil && len(routeMap[route.ID]) > 0 {
		route.Children = routeMap[route.ID]
		for i := 0; i < len(route.Children); i++ {
			rs.GetChildrenList(&route.Children[i], routeMap)
		}
	} else {
		return
	}
}
