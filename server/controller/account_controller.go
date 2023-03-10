package controller

import (
	"github.com/csyourui/wechat_server/pkg/ginserv"
	"github.com/csyourui/wechat_server/pkg/log"
	"github.com/csyourui/wechat_server/pkg/utils"
	"github.com/csyourui/wechat_server/server/service"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type AccountController struct {
	ex *service.OfficialAccount
}

func NewAccountController(ex *service.OfficialAccount) *AccountController {
	return &AccountController{ex}
}

func (ctrl *AccountController) Serve(c *gin.Context) (ginserv.Result, error) {
	log.Logger.Debug("Request Serve", c.Request)
	// 传入request和responseWriter
	ex := ctrl.ex
	server := ex.Oo.GetServer(c.Request, c.Writer)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//TODO
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}

		//article1 := message.NewArticle("测试图文1", "图文描述", "", "")
		//articles := []*message.Article{article1}
		//news := message.NewNews(articles)
		//return &message.Reply{MsgType: message.MsgTypeNews, MsgData: news}

		//voice := message.NewVoice(mediaID)
		//return &message.Reply{MsgType: message.MsgTypeVoice, MsgData: voice}

		//
		//video := message.NewVideo(mediaID, "标题", "描述")
		//return &message.Reply{MsgType: message.MsgTypeVideo, MsgData: video}

		//music := message.NewMusic("标题", "描述", "音乐链接", "HQMusicUrl", "缩略图的媒体id")
		//return &message.Reply{MsgType: message.MsgTypeMusic, MsgData: music}

		//多客服消息转发
		//transferCustomer := message.NewTransferCustomer("")
		//return &message.Reply{MsgType: message.MsgTypeTransfer, MsgData: transferCustomer}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Logger.Error("Serve Error, err=", err)
		return nil, nil
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		log.Logger.Error("Serve Error, err=", err)
		return nil, nil
	}
	return nil, nil
}

// GetAccessToken 获取ak
func (ctrl *AccountController) GetAccessToken(c *gin.Context) (ginserv.Result, error) {
	ex := ctrl.ex
	ak, err := ex.Oo.GetAccessToken()
	if err != nil {
		log.Logger.Error("get ak error, err=", err)
		utils.RenderError(c, err)
		return nil, nil
	}
	utils.RenderSuccess(c, ak)
	return nil, nil
}

// GetCallbackIP ...
func (ctrl *AccountController) GetCallbackIP(c *gin.Context) (ginserv.Result, error) {
	ex := ctrl.ex
	ipList, err := ex.Oo.GetBasic().GetCallbackIP()
	if err != nil {
		log.Logger.Error("GetCallbackIP error, err=", err)
		utils.RenderError(c, err)
		return nil, nil
	}
	utils.RenderSuccess(c, ipList)
	return nil, nil
}

// GetAPIDomainIP ...
func (ctrl *AccountController) GetAPIDomainIP(c *gin.Context) (ginserv.Result, error) {
	ex := ctrl.ex
	ipList, err := ex.Oo.GetBasic().GetAPIDomainIP()
	if err != nil {
		log.Logger.Error("GetAPIDomainIP error, err=", err)
		utils.RenderError(c, err)
		return nil, nil
	}
	utils.RenderSuccess(c, ipList)
	return nil, nil
}

// ClearQuota
func (ctrl *AccountController) ClearQuota(c *gin.Context) (ginserv.Result, error) {
	ex := ctrl.ex
	err := ex.Oo.GetBasic().ClearQuota()
	if err != nil {
		log.Logger.Error("ClearQuota error, err=", err)
		utils.RenderError(c, err)
		return nil, nil
	}
	utils.RenderSuccess(c, "success")
	return nil, nil
}
