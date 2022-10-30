package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/suyuan32/simple-utils/middleware"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/utils"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"

	"github.com/suyuan32/simple-admin-file/api/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	DB        *gorm.DB
	Redis     *redis.Redis
	Casbin    *casbin.SyncedEnforcer
	Authority rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := c.DatabaseConf.NewGORM()
	logx.Must(err)

	rds := c.RedisConf.NewRedis()

	// initialize casbin
	cbn := utils.NewCasbin(db)

	return &ServiceContext{
		Config:    c,
		DB:        db,
		Redis:     rds,
		Casbin:    cbn,
		Authority: middleware.NewAuthorityMiddleware(cbn, rds).Handle,
	}
}
