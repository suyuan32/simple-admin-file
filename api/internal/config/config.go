package config

import (
	"github.com/zeromicro/go-zero/core/stores/gormsql"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth         rest.AuthConf
	UploadConf   UploadConf
	DatabaseConf gormsql.GORMConf
	RedisConf    redis.RedisConf
}

type UploadConf struct {
	MaxImageSize     int64
	MaxVideoSize     int64
	MaxAudioSize     int64
	MaxOtherSize     int64
	PrivateStorePath string
	PublicStorePath  string
}
