package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-core/rpc/coreclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/suyuan32/simple-admin-file/internal/utils/cloud"

	"github.com/suyuan32/simple-admin-file/ent"
	"github.com/suyuan32/simple-admin-file/internal/config"
	i18n2 "github.com/suyuan32/simple-admin-file/internal/i18n"
	"github.com/suyuan32/simple-admin-file/internal/middleware"
)

type ServiceContext struct {
	Config       config.Config
	DB           *ent.Client
	Casbin       *casbin.Enforcer
	Authority    rest.Middleware
	Trans        *i18n.Translator
	CoreRpc      coreclient.Core
	CloudStorage *cloud.CloudServiceGroup
}

func NewServiceContext(c config.Config) *ServiceContext {
	entOpts := []ent.Option{
		ent.Log(logx.Info),
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
	}

	if c.DatabaseConf.Debug {
		entOpts = append(entOpts, ent.Debug())
	}

	db := ent.NewClient(entOpts...)

	rds := c.RedisConf.MustNewUniversalRedis()

	cbn := c.CasbinConf.MustNewCasbinWithOriginalRedisWatcher(c.CasbinDatabaseConf.Type,
		c.CasbinDatabaseConf.GetDSN(), c.RedisConf)

	trans := i18n.NewTranslator(c.I18nConf, i18n2.LocaleFS)

	return &ServiceContext{
		Config:       c,
		DB:           db,
		Casbin:       cbn,
		CoreRpc:      coreclient.NewCore(zrpc.MustNewClient(c.CoreRpc)),
		Authority:    middleware.NewAuthorityMiddleware(cbn, rds, trans).Handle,
		Trans:        trans,
		CloudStorage: cloud.NewCloudServiceGroup(db),
	}
}
