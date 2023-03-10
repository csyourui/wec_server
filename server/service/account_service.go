package service

import (
	"github.com/csyourui/wechat_server/pkg/log"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/spf13/viper"
)

type OfficialAccount struct {
	wc *wechat.Wechat
	Oo *officialaccount.OfficialAccount
}

func NewOfficialAccount(conf *viper.Viper, wc *wechat.Wechat) *OfficialAccount {
	//init config
	offCfg := &offConfig.Config{
		AppID:          conf.GetString("officialAccountConfig.appId"),
		AppSecret:      conf.GetString("officialAccountConfig.appSecret"),
		Token:          conf.GetString("officialAccountConfig.token"),
		EncodingAESKey: conf.GetString("officialAccountConfig.encodingAESKey"),
	}
	log.Logger.Debug("offCfg=%+v", offCfg)
	officialAccount := wc.GetOfficialAccount(offCfg)
	return &OfficialAccount{
		wc: wc,
		Oo: officialAccount,
	}
}
