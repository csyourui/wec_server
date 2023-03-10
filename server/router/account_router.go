package router

import (
	"github.com/csyourui/wechat_server/pkg/ginserv"
	"github.com/csyourui/wechat_server/pkg/ginserv/route"
	"github.com/csyourui/wechat_server/server"
	"github.com/csyourui/wechat_server/server/controller"
)

func RouteAccountCtrl(root ginserv.RouterGroup, ctrl *controller.AccountController) {
	g := root.Group("/v1/")
	routes := []*route.Route{
		route.New(g.Any, "/serve", ctrl.Serve),
		route.New(g.GET, "/oa/basic/get_access_token", ctrl.GetAccessToken),
		route.New(g.GET, "/oa/basic/get_callback_ip", ctrl.GetCallbackIP),
		route.New(g.GET, "/oa/basic/get_api_domain_ip", ctrl.GetAPIDomainIP),
		route.New(g.GET, "/oa/basic/clear_quota", ctrl.ClearQuota),
	}
	route.Bind(routes, server.ErrorCode)
}
