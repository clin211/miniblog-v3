# MiniBlog - Go-Zero 微服务博客系统

基于 Go-Zero 框架构建的轻量级博客系统，采用微服务架构设计，支持 Docker 容器化开发环境和热重载。

## 技术栈

- **框架**: [Go-Zero](https://go-zero.dev/)
- **数据库**: MySQL
- **缓存**: Redis
- **消息队列**: Kafka
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
│       ├── docker-compose.yml     # 应用服务配置
│       ├── docker-compose.env.yml # 基础环境服务配置（MySQL、Redis、Kafka）
│       ├── Dockerfile              # 开发环境 Dockerfile
│       ├── air/                    # Air 热重载配置
│       └── nginx/                  # Nginx 配置
├── tmp/                       # 临时文件目录（Air 热重载使用）
├── go.mod                     # Go 模块文件
├── go.sum                     # Go 模块依赖
├── Makefile                   # 构建和开发命令
└── request.http               # API 测试用例
```

## 快速开始

### 方式一：完整 Docker 开发环境（推荐）

#### 1. 启动基础环境服务（必须先启动）

⚠️ **重要**: 必须先启动基础环境服务，创建共享网络和依赖服务：

```bash
# 启动基础环境服务（MySQL、Redis、Kafka、可观测性组件）
make env-start

# 等待所有基础服务启动完成（大约 30-60 秒）
docker compose -f docker-compose.env.yml ps
```

#### 2. 启动应用服务

基础环境启动完成后，再启动 Go 应用服务：

```bash
# 后台启动应用开发环境（推荐）
make dev

# 或前台启动（可看实时日志）
make dev-fg

# 查看服务状态
make dev-status

# 查看日志
make dev-logs

# 停止应用环境
make dev-stop
```

#### 3. 完整的启动流程

```bash
# 完整启动流程
cd miniblog-v3/deploy/dev

# 1. 启动基础环境（网络、数据库、消息队列、可观测性）
make env-start

# 2. 等待基础服务就绪
sleep 30

# 3. 启动应用服务
make dev

# 4. 验证所有服务都在同一网络中
docker network ls
docker network inspect miniblog-network
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

### 🌐 网络架构

所有服务运行在同一个 Docker 网络 `miniblog-network` 中：

```
┌──────────────────────────────────────────────────────────────┐
│                    miniblog-network                          │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                 Application Layer                        │ │
│  │                                                         │ │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │ │
│  │  │   Nginx     │    │  User API   │    │  User RPC   │  │ │
│  │  │   Gateway   │    │   Service   │    │   Service   │  │ │
│  │  │   (8099)    │    │   (8888)    │    │   (8889)    │  │ │
│  │  └─────────────┘    └─────────────┘    └─────────────┘  │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                Infrastructure Layer                      │ │
│  │                                                         │ │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────┐ │ │
│  │  │  MySQL  │  │  Redis  │  │  Kafka  │  │ Zookeeper   │ │ │
│  │  │ (3306)  │  │ (6379)  │  │ (9092)  │  │   (2181)    │ │ │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                Observability Layer                       │ │
│  │                                                         │ │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌───────┐ │ │
│  │  │Prometheus  │ │  Grafana   │ │ Elasticsearch│ │Kibana │ │ │
│  │  │  (19090)   │ │  (3000)    │ │   (9200)     │ │(5601) │ │ │
│  │  └────────────┘ └────────────┘ └────────────┘ └───────┘ │ │
│  │                                                         │ │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐          │ │
│  │  │   Jaeger   │ │  Filebeat  │ │  go-stash  │          │ │
│  │  │  (16686)   │ │     -      │ │     -      │          │ │
│  │  └────────────┘ └────────────┘ └────────────┘          │ │
│  └─────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

### 🔄 启动依赖关系

**重要**: 必须按以下顺序启动，确保网络和依赖服务就绪：

1. **基础环境** (`docker-compose.env.yml`) → 创建网络和基础服务
2. **应用服务** (`docker-compose.yml`) → 加入已存在的网络

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

#### 基础环境服务（MySQL、Redis、Kafka）

| 命令 | 说明 |
|------|------|
| `make env-start` | 启动基础环境服务（自动创建网络） |
| `make env-stop` | 停止基础环境服务 |
| `make env-clean` | 清理基础环境和网络 |

#### 应用服务（Go 微服务）

| 命令 | 说明 |
|------|------|
| `make dev` | 后台启动应用开发环境（推荐） |
| `make dev-fg` | 前台启动应用开发环境（看实时日志） |
| `make dev-stop` | 停止应用开发环境 |
| `make dev-restart` | 重启应用开发环境 |
| `make dev-clean` | 清理应用环境（容器、镜像、数据卷） |
| `make dev-logs` | 查看应用服务日志 |
| `make dev-status` | 查看应用服务状态 |

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

### 基础服务已配置

项目已经配置了完整的基础服务和可观测性套件：

#### 基础服务

- **MySQL 8.0**: 数据库服务（端口 3306）
- **Redis 7**: 缓存服务（端口 6379）
- **Kafka 7.4.0**: 消息队列服务（端口 9092）

#### 可观测性套件

- **📊 指标监控**: Prometheus（19090）+ Grafana（3000）
- **📝 日志分析**: Elasticsearch（9200）+ Kibana（5601）+ Filebeat + go-stash
- **🔍 链路追踪**: Jaeger（16686）
- **🎛️ 管理界面**: Kafka UI（8080）

详细配置说明请参考：[环境搭建指南](docs/01%20环境搭建.md)

### 连接基础服务

在应用代码中连接这些服务：

```yaml
# 应用配置示例
database:
  host: mysql
  port: 3306
  username: miniblog_user
  password: j478EaZGDNPUbnXb

redis:
  host: redis  
  port: 6379
  password: weZ2014P89rlTuWe

kafka:
  brokers:
    - kafka:29092

# 可观测性配置
telemetry:
  # 指标收集
  metrics:
    enabled: true
    endpoint: http://prometheus:19090
  
  # 链路追踪
  tracing:
    enabled: true
    endpoint: http://jaeger:14268
    service_name: miniblog
    
  # 日志配置
  logging:
    level: info
    output: /var/log/app/app.log
    format: json
```

### 快速访问链接

| 服务 | 访问地址 | 用户名 | 密码 | 说明 |
|------|----------|--------|------|------|
| Grafana | <http://localhost:3000> | admin | admin123 | 监控可视化 |
| Kibana | <http://localhost:5601> | - | - | 日志分析 |
| Jaeger | <http://localhost:16686> | - | - | 链路追踪 |
| Kafka UI | <http://localhost:8080> | - | - | 消息队列管理 |
| Prometheus | <http://localhost:19090> | - | - | 指标查询 |

### 集成应用监控

可观测性组件已经部署完成，接下来需要在应用代码中集成：

1. **集成 Prometheus 指标**: 在 go-zero 中启用 metrics
2. **集成 Jaeger 追踪**: 配置分布式链路追踪
3. **集成结构化日志**: 输出 JSON 格式日志到指定目录
4. **配置告警规则**: 在 Grafana 中设置监控告警

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
