package router

import (
	"github.com/csyourui/wechat_server/pkg/ginserv"
	"github.com/csyourui/wechat_server/pkg/ginserv/route"
	"github.com/csyourui/wechat_server/server"
	"github.com/csyourui/wechat_server/server/controller"
)

func RouteBaseCtrl(root ginserv.RouterGroup, ctrl *controller.BaseController) {
	g := root.Group("/base/")
	routes := []*route.Route{
		route.New(g.GET, "/v1/", ctrl.Welcome),
	}
	route.Bind(routes, server.ErrorCode)
}
