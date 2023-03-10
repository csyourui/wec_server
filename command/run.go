package command

import (
	"context"
	"github.com/csyourui/wechat_server/pkg/comm"
	"github.com/csyourui/wechat_server/pkg/config"
	"github.com/csyourui/wechat_server/pkg/ginserv"
	yredis "github.com/csyourui/wechat_server/pkg/redis"
	"github.com/csyourui/wechat_server/pkg/serv"
	"github.com/csyourui/wechat_server/server"
	"github.com/csyourui/wechat_server/server/controller"
	"github.com/csyourui/wechat_server/server/router"
	"github.com/csyourui/wechat_server/server/service"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func run(conf *viper.Viper) {
	newConfig := func() (*viper.Viper, error) {
		err := config.LoadConfig(conf, "", map[string]interface{}{})
		return conf, err
	}

	app := fx.New(
		fx.Provide(
			newConfig,
			InitWechat,
			service.NewOfficialAccount,
			controller.NewAccountController,

			serv.New,
			yredis.NewRedisClusterClient,
			controller.NewBaseController,
			ginserv.NewEngine,
			ginserv.NewServer,
			ginserv.NewAPIGroup,
		),
		fx.Invoke(
			ginserv.LoadGlobalMiddlewares,
			router.RouteBaseCtrl,
			router.RouteAccountCtrl,
			server.Server,
		),
	)
	defer func() {
		if err := app.Stop(context.Background()); err != nil {
			panic(err)
		}
	}()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

//InitWechat 获取wechat实例
//在这里已经设置了全局cache，则在具体获取公众号/小程序等操作实例之后无需再设置，设置即覆盖

func InitWechat(conf *viper.Viper) *wechat.Wechat {
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:        conf.GetString("redis.host"),
		Password:    conf.GetString("redis.password"),
		Database:    conf.GetInt("redis.db"),
		MaxActive:   conf.GetInt("redis.maxActive"),
		MaxIdle:     conf.GetInt("redis.maxIdle"),
		IdleTimeout: conf.GetInt("redis.idleTimeout"),
	}
	ctx := context.Background()
	redisCache := cache.NewRedis(ctx, redisOpts)
	wc.SetCache(redisCache)
	return wc
}

func init() {
	Root.AddCommand(comm.NewRunCommand("website", "run", run))
}
