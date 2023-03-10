package controller

import (
	"github.com/csyourui/wechat_server/pkg/ginserv"
	"github.com/csyourui/wechat_server/pkg/log"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (ctrl *BaseController) Welcome(c *gin.Context) (result ginserv.Result, err error) {
	log.Logger.Debug("Welcome")
	result = ginserv.Result{
		"Code":    "200",
		"Message": "Welcome to wechat_server",
	}
	return result, nil
}
