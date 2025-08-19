# grpcurl 使用指南

## 什么是 grpcurl？

grpcurl 是一个命令行工具，类似于 curl，但专门用于 gRPC 服务。它允许你通过命令行直接调用 gRPC 方法，无需编写客户端代码。

### 设计理念

grpcurl 的设计遵循以下原则：

1. **类似 curl 的体验**：如果你熟悉 curl，那么使用 grpcurl 会感觉很自然
2. **反射优先**：优先使用 gRPC 反射 API 来获取服务信息
3. **JSON 友好**：使用 JSON 格式进行数据交换，便于调试
4. **跨平台**：支持 Windows、macOS 和 Linux

### 核心功能

- **服务发现**：列出所有可用的服务和方法
- **方法调用**：直接调用 gRPC 方法
- **类型检查**：验证请求和响应的数据结构
- **错误处理**：显示详细的错误信息
- **格式支持**：支持 JSON 和文本输出格式

### 使用场景

1. **开发和调试**：快速测试 gRPC 服务
2. **API 探索**：了解服务的可用方法
3. **自动化测试**：在 CI/CD 流程中测试服务
4. **故障排查**：诊断生产环境中的问题
5. **文档生成**：生成 API 文档

## 安装 grpcurl

```bash
# 使用 Go 安装
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 验证安装
grpcurl --version
```

## 工作原理

### 1. 反射 API

grpcurl 使用 gRPC 反射 API 来获取服务信息：

```go
// 服务器端需要启用反射
import "google.golang.org/grpc/reflection"

func main() {
    s := grpc.NewServer()
    // 注册你的服务
    pb.RegisterUserServiceServer(s, &server{})
    
    // 启用反射服务（grpcurl 需要）
    reflection.Register(s)
    
    s.Serve(lis)
}
```

### 2. 协议缓冲区

grpcurl 理解 Protocol Buffers 定义：

```protobuf
service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
}

message GetUserRequest {
    string id = 1;
}

message GetUserResponse {
    string id = 1;
    string name = 2;
    string email = 3;
}
```

### 3. JSON 序列化

grpcurl 将 JSON 转换为 Protocol Buffers：

```bash
# JSON 输入
grpcurl -d '{"id": "123"}' localhost:50051 user.UserService.GetUser

# 转换为 Protocol Buffers
GetUserRequest{Id: "123"}
```

## 完整的 gRPC 服务器代码

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/reflection"  // 添加这行
    "google.golang.org/grpc/status"

    pb "github.com/clin211/miniblog-v3/examples/grpc_service/pb"
)

// ... 服务实现代码 ...

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &server{})
    
    // 启用反射服务（grpcurl 需要）
    reflection.Register(s)  // 添加这行

    log.Printf("gRPC 服务器启动在端口 :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

## 使用 grpcurl 测试

### 1. 启动服务

```bash
cd examples/grpc_service
go run main.go pb/user.pb.go pb/user_grpc.pb.go
```

### 2. 基本测试命令

```bash
# 列出所有服务
grpcurl -plaintext localhost:50051 list

# 列出特定服务的方法
grpcurl -plaintext localhost:50051 list user.UserService

# 查看方法详情
grpcurl -plaintext localhost:50051 describe user.UserService.GetUser

# 调用方法
grpcurl -plaintext -d '{"id": "123"}' localhost:50051 user.UserService.GetUser
```

### 3. 完整测试脚本

运行 `test_working.sh` 脚本进行完整测试：

```bash
./test_working.sh
```

## grpcurl 常用参数

- `-plaintext`: 使用明文连接（不加密）
- `-d`: 指定请求数据（JSON 格式）
- `-v`: 详细输出
- `-format`: 输出格式（json, text）

## 错误处理

gRPC 使用标准状态码：

- `OK (0)`: 成功
- `INVALID_ARGUMENT (3)`: 无效参数
- `NOT_FOUND (5)`: 资源未找到
- `INTERNAL (13)`: 内部错误

## 调试技巧

1. **使用 -v 参数查看详细信息**
2. **检查服务是否启动**：

   ```bash
   lsof -i :50051
   ```

3. **使用 grpcui 进行交互式测试**：

   ```bash
   go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
   grpcui -plaintext localhost:50051
   ```

## 注意事项

1. **反射服务**：grpcurl 需要反射 API 才能工作
2. **端口检查**：确保服务在正确的端口启动
3. **协议版本**：使用 `-plaintext` 进行非加密连接
4. **JSON 格式**：请求数据必须是有效的 JSON

## grpcurl 与其他工具对比

### grpcurl vs curl

| 特性 | grpcurl | curl |
|------|---------|------|
| 协议支持 | gRPC (HTTP/2) | HTTP/1.1, HTTP/2 |
| 数据格式 | Protocol Buffers | JSON, XML, 文本 |
| 类型安全 | 强类型 | 弱类型 |
| 服务发现 | 通过反射 API | 需要 API 文档 |
| 性能 | 更高（二进制） | 较低（文本） |

### grpcurl vs grpcui

| 特性 | grpcurl | grpcui |
|------|---------|--------|
| 界面 | 命令行 | Web UI |
| 交互性 | 单次调用 | 交互式 |
| 脚本化 | 适合自动化 | 适合手动测试 |
| 学习曲线 | 简单 | 中等 |
| 部署 | 轻量级 | 需要 Web 服务器 |

### grpcurl vs Postman

| 特性 | grpcurl | Postman |
|------|---------|---------|
| 协议支持 | gRPC | HTTP REST |
| 平台 | 跨平台命令行 | 桌面应用 |
| 团队协作 | 代码版本控制 | 内置协作功能 |
| 自动化 | 脚本友好 | 内置测试功能 |
| 学习成本 | 低 | 中等 |

## 最佳实践

### 1. 开发环境

```bash
# 启用反射服务（仅开发环境）
reflection.Register(s)

# 使用详细输出进行调试
grpcurl -v -plaintext localhost:50051 list
```

### 2. 生产环境

```bash
# 禁用反射服务以提高安全性
# reflection.Register(s)  // 注释掉这行

# 使用 TLS 加密
grpcurl -cert client.crt -key client.key localhost:50051 list
```

### 3. 自动化测试

```bash
#!/bin/bash
# 测试脚本示例
grpcurl -plaintext -d '{"id": "123"}' localhost:50051 user.UserService.GetUser | jq '.id' | grep -q "123" && echo "测试通过" || echo "测试失败"
```

### 4. 错误处理

```bash
# 检查服务状态
grpcurl -plaintext localhost:50051 list || echo "服务不可用"

# 处理错误响应
grpcurl -plaintext -d '{"id": "0"}' localhost:50051 user.UserService.GetUser 2>&1 | grep "InvalidArgument"
```

## 常见问题

### Q: 为什么需要反射 API？

A: grpcurl 需要了解服务的结构才能正确调用方法。反射 API 提供了这种信息。

### Q: 生产环境可以启用反射吗？

A: 不建议。反射 API 会暴露服务结构信息，可能带来安全风险。

### Q: 如何在没有反射的情况下使用 grpcurl？

A: 可以使用 `-proto` 参数指定 .proto 文件：

```bash
grpcurl -proto user.proto -plaintext localhost:50051 user.UserService.GetUser
```

### Q: grpcurl 支持流式调用吗？

A: 是的，支持 unary、server streaming、client streaming 和 bidirectional streaming。
