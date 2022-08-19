package svc

import (
	"github.com/suyuan32/simple-admin-file/api/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"log"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := c.DB.NewGORM()
	if err != nil {
		log.Fatal(err.Error())
	}
	rds := c.RedisConf.NewRedis()
	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rds,
	}
}
