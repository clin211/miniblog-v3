// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package token

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Config 包括 token 包的配置选项.
type Config struct {
	// Secret 用于签发和解析 token 的密钥.
	Secret string
	// IdentityKey 是 token 中用户身份的键.
	IdentityKey string
	// Expiration 是签发的 token 过期时间
	Expiration time.Duration
	// Issuer 是 token 的签发者
	Issuer string
	// Audience 是 token 的目标受众
	Audience string
}

// Claims 定义 JWT 的声明结构
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	defaultConfig = Config{
		Secret:      "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		IdentityKey: "user_id",
		Expiration:  24 * time.Hour,
		Issuer:      "miniblog",
		Audience:    "miniblog_users",
	}
	globalConfig Config
	once         sync.Once
)

// Init 初始化 token 包配置
func Init(config Config) {
	once.Do(func() {
		if config.Secret == "" {
			config.Secret = defaultConfig.Secret
		}
		if config.IdentityKey == "" {
			config.IdentityKey = defaultConfig.IdentityKey
		}
		if config.Expiration == 0 {
			config.Expiration = defaultConfig.Expiration
		}
		if config.Issuer == "" {
			config.Issuer = defaultConfig.Issuer
		}
		if config.Audience == "" {
			config.Audience = defaultConfig.Audience
		}

		globalConfig = config
	})
}

// getConfig 获取当前配置
func getConfig() Config {
	if globalConfig.Secret == "" {
		Init(defaultConfig)
	}
	return globalConfig
}

// Sign 签发 JWT Token
func Sign(userID string) (string, time.Time, error) {
	config := getConfig()
	now := time.Now()
	expireAt := now.Add(config.Expiration)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Audience:  []string{config.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("签发 token 失败: %w", err)
	}

	return tokenString, expireAt, nil
}

// Parse 解析 JWT Token
func Parse(tokenString string) (*Claims, error) {
	config := getConfig()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保 token 加密算法是预期的加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析 token 失败: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// ParseRequest 从请求头中获取令牌，并将其传递给 Parse 函数以解析令牌.
// 支持 HTTP 和 gRPC 两种上下文类型
func ParseRequest(ctx interface{}) (*Claims, error) {
	var (
		token string
		err   error
	)

	switch typed := ctx.(type) {
	// 使用标准 HTTP 请求
	case *http.Request:
		header := typed.Header.Get("Authorization")
		if len(header) == 0 {
			return nil, errors.New("缺少 Authorization 头")
		}

		// 从请求头中取出 token
		_, _ = fmt.Sscanf(header, "Bearer %s", &token) // 解析 Bearer token
		if token == "" {
			return nil, errors.New("token 为空")
		}
	// 使用 google.golang.org/grpc 框架开发的 gRPC 服务
	default:
		if ctxTyped, ok := ctx.(context.Context); ok {
			token, err = auth.AuthFromMD(ctxTyped, "Bearer")
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
			}
		} else {
			return nil, errors.New("不支持的上下文类型")
		}
	}

	return Parse(token)
}
