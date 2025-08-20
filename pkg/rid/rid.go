// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package rid

import (
	"context"

	"github.com/clin211/miniblog-v3/pkg/id"
)

const defaultABC = "abcdefghijklmnopqrstuvwxyz1234567890"

type ResourceID string

const (
	// UserID 定义用户资源标识符，前缀为 miniblog-user 的缩写 mu
	UserID ResourceID = "mu"
)

// String 将资源标识符转换为字符串.
func (rid ResourceID) String() string {
	return string(rid)
}

// New 创建带前缀的唯一标识符.
func (rid ResourceID) New() string {
	sf := id.NewSonyflake()
	// 通过 Sonyflake 生成数值 ID，再编码为短码
	num := sf.Id(context.Background())

	// 使用自定义选项生成唯一标识符
	uniqueStr := id.NewCode(
		num,
		id.WithCodeChars([]rune(defaultABC)),
		id.WithCodeL(6),
		id.WithCodeSalt(Salt()),
	)
	return rid.String() + "-" + uniqueStr
}
