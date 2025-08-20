// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	// 初始化 token 包
	config := Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	Init(config)

	// 测试签发 token
	tokenString, expireAt, err := Sign("user_123")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	assert.True(t, expireAt.After(time.Now()))

	// 验证 token 格式
	assert.Contains(t, tokenString, ".")
}

func TestParse(t *testing.T) {
	// 初始化 token 包
	config := Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	Init(config)

	// 签发 token
	tokenString, _, err := Sign("user_123")
	assert.NoError(t, err)

	// 解析 token
	claims, err := Parse(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user_123", claims.UserID)
}

func TestParseInvalidToken(t *testing.T) {
	// 初始化 token 包
	config := Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	Init(config)

	// 测试解析无效 token
	_, err := Parse("invalid-token")
	assert.Error(t, err)
}

func TestParseRequest(t *testing.T) {
	// 初始化 token 包
	config := Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	Init(config)

	// 签发 token
	tokenString, _, err := Sign("user_123")
	assert.NoError(t, err)

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// 解析 token
	claims, err := ParseRequest(req)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user_123", claims.UserID)
}

func TestParseRequestNoHeader(t *testing.T) {
	// 初始化 token 包
	config := Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	Init(config)

	// 创建测试请求（无 Authorization 头）
	req := httptest.NewRequest("GET", "/test", nil)

	// 解析 token
	_, err := ParseRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "缺少 Authorization 头")
}

func TestParseRequestInvalidFormat(t *testing.T) {
	// 初始化 token 包
	config := Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	Init(config)

	// 创建测试请求（格式错误的 Authorization 头）
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")

	// 解析 token
	_, err := ParseRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token 为空")
}

func TestInit(t *testing.T) {
	// 测试初始化
	config := Config{
		Secret:      "custom-secret",
		IdentityKey: "custom_user_id",
		Expiration:  2 * time.Hour,
		Issuer:      "custom-issuer",
		Audience:    "custom-audience",
	}
	Init(config)

	// 测试签发 token
	tokenString, _, err := Sign("user_456")
	assert.NoError(t, err)

	// 解析 token
	claims, err := Parse(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, "user_456", claims.UserID)
}

func TestTokenExpiration(t *testing.T) {
	// 初始化 token 包
	config := Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	Init(config)

	// 签发 token
	tokenString, _, err := Sign("user_123")
	assert.NoError(t, err)

	// 解析 token 应该成功
	claims, err := Parse(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user_123", claims.UserID)
}

func TestDefaultConfig(t *testing.T) {
	// 测试默认配置
	Init(Config{}) // 使用空配置，应该使用默认值

	// 测试签发 token
	tokenString, _, err := Sign("user_789")
	assert.NoError(t, err)

	// 解析 token
	claims, err := Parse(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, "user_789", claims.UserID)
}
