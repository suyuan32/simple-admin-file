package config

import (
	"github.com/suyuan32/simple-admin-tools/plugins/registry/consul"
	"github.com/zeromicro/go-zero/core/stores/gormsql"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf `yaml:",inline"`
	Auth          Auth             `json:"auth" yaml:"Auth"`
	UploadConf    UploadConf       `json:"UploadConf" yaml:"UploadConf"`
	DB            gormsql.GORMConf `json:"DatabaseConf" yaml:"DatabaseConf"`
	RedisConf     redis.RedisConf  `json:"RedisConf" yaml:"RedisConf"`
}

type UploadConf struct {
	MaxImageSize     int64  `json:"MaxImageSize" yaml:"MaxImageSize"`
	MaxVideoSize     int64  `json:"MaxVideoSize" yaml:"MaxVideoSize"`
	MaxAudioSize     int64  `json:"MaxAudioSize" yaml:"MaxAudioSize"`
	MaxOtherSize     int64  `json:"MaxOtherSize" yaml:"MaxOtherSize"`
	PrivateStorePath string `json:"PrivateStorePath" yaml:"PrivateStorePath"`
	PublicStorePath  string `json:"PublicStorePath" yaml:"PublicStorePath"`
}

type Auth struct {
	AccessSecret string `json:"accessSecret" yaml:"AccessSecret"`
	AccessExpire int64  `json:"accessExpire" yaml:"AccessExpire"`
}

type ConsulConfig struct {
	Consul consul.Conf
}
