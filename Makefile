# Go项目构建工具
.PHONY: fmt lint test build dev

# 项目根目录
ROOT_DIR := $(shell pwd)

# 服务配置
# 默认包含所有发现的服务
# 要排除的服务添加到DISABLED_SERVICES变量

# 动态发现所有服务
ALL_SERVICES := $(wildcard apps/*/api) $(wildcard apps/*/rpc)

# 要排除的服务(设置为空表示包含所有服务)
# 示例: DISABLED_SERVICES = apps/user/api apps/user/rpc
DISABLED_SERVICES ?=

# 激活的服务列表
ENABLED_SERVICES = $(filter-out $(DISABLED_SERVICES),$(ALL_SERVICES))

# 格式化代码
fmt:
	go fmt ./...
	goimports -w -l .

# 静态代码检查
lint:
	golangci-lint run ./...

# 运行测试
test:
	go test -v -cover ./...

# 构建所有服务
build:
	@for service in $(ENABLED_SERVICES); do \
		service_name=$$(basename $$(dirname $$service))-$$(basename $$service); \
		echo "Building $$service_name..."; \
		go build -o $(ROOT_DIR)/bin/$$service_name $(ROOT_DIR)/$$service; \
	done

# 开发模式(启动所有启用服务)
dev:
	@for service in $(ENABLED_SERVICES); do \
		service_name=$$(basename $$(dirname $$service))-$$(basename $$service); \
		echo "Starting $$service_name..."; \
		cd $(ROOT_DIR)/$$(dirname $$service) && \
		if [ "$$(basename $$service)" = "api" ]; then \
			CONFIG_FILE=$(ROOT_DIR)/$$service/etc/api-api.yaml; \
		else \
			CONFIG_FILE=$(ROOT_DIR)/$$service/etc/rpc-api.yaml; \
		fi && \
		air -c $(ROOT_DIR)/.air.toml -build.cmd "go build -o $(ROOT_DIR)/tmp/$$service_name $(ROOT_DIR)/$$service" -build.bin "$(ROOT_DIR)/tmp/$$service_name" -- -f $$CONFIG_FILE & \
	done; \
	echo "Services started in background. Use 'pkill air' to stop."

# 安装开发工具链
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/air-verse/air@latest
