# MiniBlog - Go-Zero 微服务博客系统

基于 Go-Zero 框架构建的轻量级博客系统，采用微服务架构设计。

## 技术栈

- **框架**: [Go-Zero](https://go-zero.dev/)
- **数据库**: MySQL (待配置)
- **缓存**: Redis (待配置)
- **API 文档**: 自动生成
- **容器化**：Docker
- **测试工具**: REST Client (VSCode 插件)

## 项目结构

```tree
miniblog-v3/
├── apps/
│   ├── blog/       # 博客服务
│   └── user/       # 用户服务
│       ├── api/    # HTTP API 服务
│       └── rpc/    # RPC 服务
├── go.mod          # Go 模块文件
└── request.http    # API 测试用例
```

## 快速开始

### 环境准备

1. 安装 Go 1.16+
2. 安装 goctl 工具:

   ```bash
   go install github.com/zeromicro/go-zero/tools/goctl@latest
   ```

### 启动服务

```bash
# 启动用户API服务
cd apps/user/api
go run api.go -f etc/api-api.yaml
```

## API 测试

项目包含 `request.http` 文件，配合 VSCode 的 [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) 插件可直接测试 API。

示例测试端点:

- `GET /from/you`
- `GET /from/me`

## 开发指南

1. 生成代码:

   ```bash
   goctl api go -api api.api -dir .
   ```

2. 修改 API 定义文件后重新生成代码

3. 添加业务逻辑到 `internal/logic/` 目录

## 贡献

欢迎提交 Pull Request。

## License

MIT
