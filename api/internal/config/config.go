package config

import (
	"github.com/suyuan32/simple-admin-common/config"
	"github.com/suyuan32/simple-admin-common/plugins/casbin"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth               rest.AuthConf
	UploadConf         UploadConf
	DatabaseConf       config.DatabaseConf
	CasbinDatabaseConf config.DatabaseConf
	RedisConf          redis.RedisConf
	CasbinConf         casbin.CasbinConf
	CoreRpc            zrpc.RpcClientConf
}

type UploadConf struct {
	MaxImageSize     int64
	MaxVideoSize     int64
	MaxAudioSize     int64
	MaxOtherSize     int64
	PrivateStorePath string
	PublicStorePath  string
}
