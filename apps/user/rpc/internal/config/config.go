package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// MySQL 数据库配置
	MySQL struct {
		Host     string `json:",default=localhost"`
		Port     int    `json:",default=3306"`
		User     string `json:",default=root"`
		Password string `json:",default=password"`
		Database string `json:",default=miniblog_user"`
		Charset  string `json:",default=utf8mb4"`
		MaxOpen  int    `json:",default=100"`
		MaxIdle  int    `json:",default=10"`
		Timeout  string `json:",default=5s"`
	}

	// Redis 缓存配置
	Redis redis.RedisConf

	// JWT 配置
	JWT struct {
		Secret      string `json:",default=your-jwt-secret-key"`
		ExpireHours int    `json:",default=24"`
	}

	// 服务配置
	Service struct {
		Name string `json:",default=user-rpc"`
	}
}
