# MiniBlog - Go-Zero 微服务博客系统

基于 Go-Zero 框架构建的轻量级博客系统，采用微服务架构设计，支持 Docker 容器化开发环境和热重载。

## 技术栈

- **框架**: [Go-Zero](https://go-zero.dev/)
- **数据库**: MySQL (待配置)
- **缓存**: Redis (待配置)
- **API 文档**: 自动生成
- **容器化**: Docker + Docker Compose
- **热重载**: [Air](https://github.com/air-verse/air)
- **网关**: Nginx
- **测试工具**: REST Client (VSCode 插件)

## 项目结构

```tree
miniblog-v3/
├── apps/                       # 微服务应用目录
│   └── user/                   # 用户服务
│       ├── api/                # HTTP API 服务
│       │   ├── etc/           # 配置文件
│       │   ├── internal/      # 内部代码
│       │   ├── user.api       # API 定义
│       │   └── user.go        # 服务入口
│       └── rpc/               # RPC 服务
│           ├── etc/           # 配置文件
│           ├── internal/      # 内部代码
│           ├── rpc.api        # RPC 定义
│           └── rpc.go         # 服务入口
├── _output/                    # 构建输出目录（make build 生成）
├── deploy/                     # 部署配置
│   └── dev/                   # 开发环境配置
│       ├── docker-compose.yml # Docker Compose 配置
│       ├── Dockerfile         # 开发环境 Dockerfile
│       ├── air/               # Air 热重载配置
│       └── nginx/             # Nginx 配置
├── tmp/                       # 临时文件目录（Air 热重载使用）
├── go.mod                     # Go 模块文件
├── go.sum                     # Go 模块依赖
├── Makefile                   # 构建和开发命令
└── request.http               # API 测试用例
```

## 快速开始

### 方式一：Docker 开发环境（推荐）

使用 Docker 和 Air 热重载的完整开发环境：

```bash
# 后台启动开发环境（推荐）
make dev

# 或前台启动（可看实时日志）
make dev-fg

# 查看服务状态
make dev-status

# 查看日志
make dev-logs

# 停止环境
make dev-stop
```

### 方式二：本地开发环境

#### 环境准备

1. 安装 Go 1.24+
2. 安装开发工具:

   ```bash
   # 安装必要工具
   make install-tools
   
   # 手动安装 goctl 工具
   go install github.com/zeromicro/go-zero/tools/goctl@latest
   ```

#### 启动服务

```bash
# 手动启动单个服务
cd apps/user/api
go run user.go -f etc/user.yaml
```

## 开发环境架构

```
┌─────────────────────────────────────────────────┐
│                   Nginx Gateway                 │
│                  (Port: 8099)                   │
│                                                 │
│  /api/user/*  ─────→  user-api:8888            │
│  /rpc/user/*  ─────→  user-rpc:8889            │
└─────────────────────────────────────────────────┘
                         │
               ┌─────────┴─────────┐
               │                   │
         ┌─────▼─────┐       ┌─────▼─────┐
         │ User API  │       │ User RPC  │
         │ Service   │       │ Service   │
         │ (8888)    │       │ (8889)    │
         └───────────┘       └───────────┘
```

### 服务列表

| 服务名称 | 容器名称 | 内部端口 | 外部端口 | 说明 |
|---------|---------|---------|---------|------|
| user-api | miniblog-user-api | 8888 | 8888 | 用户API服务 |
| user-rpc | miniblog-user-rpc | 8889 | 8889 | 用户RPC服务 |
| nginx | miniblog-nginx | 8099 | 8099 | API网关 |

### 访问地址

- **网关健康检查**: `http://localhost:8099/health`
- **用户API服务**: `http://localhost:8099/api/user/*`
- **用户RPC服务**: `http://localhost:8099/rpc/user/*`
- **直接访问API**: `http://localhost:8888`
- **直接访问RPC**: `http://localhost:8889`

## 开发工作流

### 代码热重载

项目使用 Air 工具实现热重载，当你修改代码时：

1. 修改 `apps/user/api/` 或 `apps/user/rpc/` 目录下的任何 Go 文件
2. Air 自动检测变化并重新编译
3. 服务自动重启，无需手动操作

### 添加新服务

1. 在 `apps/` 目录下创建新服务目录
2. 使用 goctl 生成服务代码
3. 修改 `deploy/dev/docker-compose.yml` 添加新服务配置
4. 在 `deploy/dev/air/` 目录下创建对应的 Air 配置文件
5. 更新 `deploy/dev/nginx/conf.d/miniblog.conf` 添加路由规则

## 常用命令

### 开发环境管理

| 命令 | 说明 |
|------|------|
| `make dev` | 后台启动开发环境（推荐） |
| `make dev-fg` | 前台启动开发环境（看实时日志） |
| `make dev-stop` | 停止开发环境 |
| `make dev-restart` | 重启开发环境 |
| `make dev-clean` | 清理环境（容器、镜像、数据卷） |
| `make dev-logs` | 查看所有服务日志 |
| `make dev-status` | 查看服务状态 |

### 代码质量工具

| 命令 | 说明 |
|------|------|
| `make build` | 构建所有服务到 _output/ 目录 |
| `make clean` | 清理构建产物 |
| `make fmt` | 格式化代码 |
| `make lint` | 静态代码检查 |
| `make test` | 运行测试 |
| `make install-tools` | 安装开发工具链 |

## API 测试

项目包含 `request.http` 文件，配合 VSCode 的 [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) 插件可直接测试 API。

### 测试端点示例

```http
### 健康检查
GET http://localhost:8099/health

### 用户API测试
GET http://localhost:8099/api/user/ping

### RPC服务测试
GET http://localhost:8099/rpc/user/ping
```

## 开发指南

### 1. 生成 API 服务代码

```bash
# 在服务目录下生成代码
cd apps/user/api
goctl api go -api user.api -dir .
```

### 2. 生成 RPC 服务代码

```bash
# 在服务目录下生成代码
cd apps/user/rpc
goctl rpc protoc rpc.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

### 3. 修改业务逻辑

- API 业务逻辑：`apps/*/api/internal/logic/`
- RPC 业务逻辑：`apps/*/rpc/internal/logic/`
- 配置文件：`apps/*/*/etc/`

### 4. 添加路由

修改对应的 `internal/handler/routes.go` 文件。

## 故障排除

### 常见问题

1. **端口冲突**

   ```bash
   # 检查端口占用
   lsof -i :8099
   lsof -i :8888  
   lsof -i :8889
   ```

2. **Docker 服务启动失败**

   ```bash
   # 查看详细日志
   cd deploy/dev && docker-compose logs [service-name]
   ```

3. **热重载不工作**
   - 检查文件是否在 Air 监听目录范围内
   - 查看 Air 构建日志排查编译错误
   - 确认 Docker 文件挂载正确

4. **网关路由问题**
   - 检查 Nginx 配置语法
   - 确认上游服务正常运行
   - 查看 Nginx 错误日志

### 重置环境

```bash
# 完全清理并重新构建
make dev-clean
make dev
```

## 性能优化

- **Go Modules 缓存**: Docker volume 缓存依赖包
- **构建缓存**: 利用 Docker 分层缓存
- **热重载优化**: Air 配置排除不必要文件监听

## 扩展功能

### 添加数据库

1. 在 `deploy/dev/docker-compose.yml` 中添加 MySQL/Redis 服务
2. 更新服务配置文件添加数据库连接信息
3. 添加数据库迁移脚本

### 添加监控

1. 集成 Prometheus + Grafana
2. 配置 go-zero 监控指标
3. 在 Nginx 中添加监控端点

## 生产部署

开发环境配置不适用于生产，生产环境需要：

1. 移除热重载工具
2. 多阶段构建优化镜像大小  
3. 生产级 Nginx 配置
4. 健康检查和监控
5. 安全配置和证书

## 贡献

欢迎提交 Pull Request 和 Issue。

## License

MIT
