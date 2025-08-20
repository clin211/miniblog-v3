# JWT 封装 Token 包设计文档

## 1. 设计思路

### 1.1 设计目标

Token 包的设计目标是提供一个轻量级、易用的 JWT 认证解决方案，专门适配 go-zero 框架，支持 HTTP 和 gRPC 两种服务类型。

### 1.2 核心特性

- **轻量级**: 只包含核心的 token 操作功能
- **幂等性**: 所有方法都是幂等的，可以安全地重复调用
- **类型安全**: 使用 Go 的类型系统确保类型安全
- **框架适配**: 专门为 go-zero 框架设计
- **双协议支持**: 同时支持 HTTP 和 gRPC 协议

### 1.3 设计原则

1. **单一职责**: 每个函数只负责一个特定的功能
2. **接口统一**: 提供统一的接口处理不同类型的请求
3. **配置灵活**: 支持灵活的配置选项
4. **错误处理**: 提供清晰的错误信息

## 2. 架构设计

### 2.1 包结构

```tree
pkg/token/
├── token.go          # 核心实现
└── token_test.go     # 单元测试
```

### 2.2 核心组件

#### 2.2.1 Config 配置结构

```go
type Config struct {
    Secret      string        // JWT 签名密钥
    IdentityKey string        // 用户身份标识键
    Expiration  time.Duration // Token 过期时间
    Issuer      string        // Token 签发者
    Audience    string        // Token 目标受众
}
```

#### 2.2.2 Claims 声明结构

```go
type Claims struct {
    UserID string `json:"user_id"` // 用户ID
    jwt.RegisteredClaims           // JWT 标准声明
}
```

### 2.3 函数设计

#### 2.3.1 初始化函数

```go
func Init(config Config)
```

- **功能**: 初始化 token 包配置
- **特性**: 使用 `sync.Once` 确保只初始化一次
- **默认值**: 提供合理的默认配置

#### 2.3.2 Token 签发函数

```go
func Sign(userID string) (string, time.Time, error)
```

- **功能**: 签发 JWT Token
- **参数**: 用户ID
- **返回**: Token 字符串、过期时间和错误

#### 2.3.3 Token 解析函数

```go
func Parse(tokenString string) (*Claims, error)
```

- **功能**: 解析 JWT Token
- **参数**: Token 字符串
- **返回**: Claims 结构体和错误

#### 2.3.4 请求解析函数

```go
func ParseRequest(ctx interface{}) (*Claims, error)
```

- **功能**: 从请求中解析 Token
- **特性**: 支持 HTTP 和 gRPC 两种上下文类型
- **实现**: 使用类型断言处理不同输入类型

## 3. 实现细节

### 3.1 配置管理

```go
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
```

- 使用全局变量存储配置
- 使用 `sync.Once` 确保线程安全的初始化
- 提供默认配置值

### 3.2 Token 签发实现

```go
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
```

### 3.3 统一请求解析

```go
func ParseRequest(ctx interface{}) (*Claims, error) {
    var (
        token string
        err   error
    )

    switch typed := ctx.(type) {
    case *http.Request:
        // HTTP 请求处理
        header := typed.Header.Get("Authorization")
        if len(header) == 0 {
            return nil, errors.New("缺少 Authorization 头")
        }
        _, _ = fmt.Sscanf(header, "Bearer %s", &token)
        if token == "" {
            return nil, errors.New("token 为空")
        }
    default:
        // gRPC 请求处理
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
```

## 4. 使用示例

### 4.1 HTTP 服务示例

#### 4.1.1 配置初始化

```go
// 在服务启动时初始化
token.Init(token.Config{
    Secret:      "your-secret-key",
    IdentityKey: "user_id",
    Expiration:  24 * time.Hour,
    Issuer:      "miniblog",
    Audience:    "miniblog_users",
})
```

#### 4.1.2 登录接口

```go
func (l *UserLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    // 验证用户名密码
    if req.Username == "admin" && req.Password == "123456" {
        // 生成 token
        tokenString, expireAt, err := token.Sign("user_123")
        if err != nil {
            return nil, err
        }

        return &types.LoginResp{
            Token:     tokenString,
            ExpiresAt: expireAt.Format(time.RFC3339),
        }, nil
    }

    return nil, errors.New("用户名或密码错误")
}
```

#### 4.1.3 受保护的接口

```go
func (l *UserLogic) GetProfile() (*types.UserProfileResp, error) {
    // 从请求中解析 token
    claims, err := token.ParseRequest(l.ctx)
    if err != nil {
        return nil, err
    }

    return &types.UserProfileResp{
        UserID: claims.UserID,
    }, nil
}
```

### 4.2 gRPC 服务示例

#### 4.2.1 服务实现

```go
func (l *UserLogic) GetProfile(in *pb.GetProfileReq) (*pb.GetProfileResp, error) {
    // 从上下文中解析 token
    claims, err := token.ParseRequest(l.ctx)
    if err != nil {
        return nil, err
    }

    return &pb.GetProfileResp{
        UserId: claims.UserID,
    }, nil
}
```

## 5. 测试策略

### 5.1 单元测试

- **Token 签发测试**: 验证 token 签发功能
- **Token 解析测试**: 验证 token 解析功能
- **请求解析测试**: 验证 HTTP 和 gRPC 请求解析
- **错误处理测试**: 验证各种错误情况
- **过期时间测试**: 验证 token 过期机制

### 5.2 测试覆盖

```go
func TestSign(t *testing.T) {
    // 测试 token 签发
    tokenString, expireAt, err := Sign("user_123")
    assert.NoError(t, err)
    assert.NotEmpty(t, tokenString)
    assert.True(t, expireAt.After(time.Now()))
}
```

func TestParseRequest(t *testing.T) {
    // 测试 HTTP 请求解析
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer "+tokenString)
    claims, err := ParseRequest(req)
    assert.NoError(t, err)
    assert.NotNil(t, claims)
    assert.Equal(t, "user_123", claims.UserID)
}

```

## 6. 最佳实践

### 6.1 配置管理

1. **密钥安全**: 使用强密钥，避免硬编码
2. **环境配置**: 根据环境调整配置参数
3. **默认值**: 提供合理的默认配置

### 6.2 错误处理

1. **错误信息**: 提供清晰的错误信息
2. **错误类型**: 区分不同类型的错误
3. **日志记录**: 记录关键操作的日志

### 6.3 性能优化

1. **缓存配置**: 避免重复读取配置
2. **连接复用**: 复用 HTTP 客户端连接
3. **内存管理**: 及时释放不需要的资源

## 7. 扩展性

### 7.1 支持的扩展

- **自定义声明**: 可以扩展 Claims 结构
- **多种算法**: 支持不同的签名算法
- **中间件集成**: 可以集成到 go-zero 中间件

### 7.2 未来规划

- **Redis 集成**: 支持 token 存储到 Redis
- **多租户**: 支持多租户场景
- **监控指标**: 添加性能监控指标

## 8. 总结

Token 包通过简洁的设计和统一的接口，为 go-zero 框架提供了完整的 JWT 认证解决方案。其特点包括：

1. **轻量级**: 只包含核心功能，易于集成
2. **类型安全**: 充分利用 Go 的类型系统
3. **双协议支持**: 同时支持 HTTP 和 gRPC
4. **配置灵活**: 支持灵活的配置选项
5. **测试完善**: 提供完整的测试覆盖

这个设计既满足了当前的需求，又为未来的扩展留下了空间。
