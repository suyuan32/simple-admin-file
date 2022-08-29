package config

import (
	"github.com/zeromicro/go-zero/core/stores/gormsql"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	UploadConf UploadConf       `json:"UploadConf" yaml:"UploadConf"`
	DB         gormsql.GORMConf `json:"DatabaseConf" yaml:"DatabaseConf"`
	RedisConf  redis.RedisConf  `json:"RedisConf" yaml:"RedisConf"`
}

type UploadConf struct {
	MaxImageSize     int64  `json:"MaxImageSize" yaml:"MaxImageSize"`
	MaxVideoSize     int64  `json:"MaxVideoSize" yaml:"MaxVideoSize"`
	MaxAudioSize     int64  `json:"MaxAudioSize" yaml:"MaxAudioSize"`
	MaxOtherSize     int64  `json:"MaxOtherSize" yaml:"MaxOtherSize"`
	PrivateStorePath string `json:"PrivateStorePath" yaml:"PrivateStorePath"`
	PublicStorePath  string `json:"PublicStorePath" yaml:"PublicStorePath"`
}
