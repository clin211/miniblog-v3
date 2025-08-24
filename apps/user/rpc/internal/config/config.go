// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

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
