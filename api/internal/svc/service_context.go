package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/suyuan32/simple-admin-core/rpc/coreclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/suyuan32/simple-admin-file/api/internal/config"
	"github.com/suyuan32/simple-admin-file/api/internal/middleware"
	"github.com/suyuan32/simple-admin-file/pkg/ent"
	i18n2 "github.com/suyuan32/simple-admin-file/pkg/i18n"
)

type ServiceContext struct {
	Config    config.Config
	DB        *ent.Client
	Redis     *redis.Redis
	Casbin    *casbin.Enforcer
	Authority rest.Middleware
	Trans     *i18n.Translator
	CoreRpc   coreclient.Core
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := ent.NewClient(
		ent.Log(logx.Info), // logger
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
		ent.Debug(), // debug mode
	)

	rds := c.RedisConf.NewRedis()
	if !rds.Ping() {
		logx.Error("initialize redis failed")
		return nil
	}

	cbn, err := c.CasbinConf.NewCasbin(c.DatabaseConf.Type, c.DatabaseConf.GetDSN())
	if err != nil {
		logx.Errorw("Initialize casbin failed", logx.Field("detail", err.Error()))
		return nil
	}

	trans := &i18n.Translator{}
	trans.NewBundle(i18n2.LocaleFS)
	trans.NewTranslator()

	return &ServiceContext{
		Config:    c,
		DB:        db,
		Redis:     rds,
		Casbin:    cbn,
		CoreRpc:   coreclient.NewCore(zrpc.MustNewClient(c.CoreRpc)),
		Authority: middleware.NewAuthorityMiddleware(cbn, rds).Handle,
		Trans:     trans,
	}
}
