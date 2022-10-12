package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"

	"github.com/suyuan32/simple-admin-file/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := c.DatabaseConf.NewGORM()
	logx.Must(err)

	rds := c.RedisConf.NewRedis()
	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rds,
	}
}
