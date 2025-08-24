// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package models

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CasbinRuleModel = (*customCasbinRuleModel)(nil)

type (
	// CasbinRuleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCasbinRuleModel.
	CasbinRuleModel interface {
		casbinRuleModel
	}

	customCasbinRuleModel struct {
		*defaultCasbinRuleModel
	}
)

// NewCasbinRuleModel returns a model for the database table.
func NewCasbinRuleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CasbinRuleModel {
	return &customCasbinRuleModel{
		defaultCasbinRuleModel: newCasbinRuleModel(conn, c, opts...),
	}
}
