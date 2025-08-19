package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Cache cache.CacheConf

	// MySQL 数据库
	Mysql struct {
		DataSource string
	}

	// JWT 配置
	JWT struct {
		Secret      string
		ExpireHours int
	}

	// 服务配置
	Service struct {
		Name string
	}
}
