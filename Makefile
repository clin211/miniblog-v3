# Go项目构建工具
.PHONY: fmt lint test build clean install-tools dev dev-fg dev-stop dev-clean dev-logs dev-restart dev-status env-start env-stop env-clean add-copyright

# 项目根目录
ROOT_DIR := $(shell pwd)
OUTPUT_DIR := $(ROOT_DIR)/_output

# 代码质量工具
fmt:
	go fmt ./...
	goimports -w -l .

lint:
	golangci-lint run ./...

test:
	go test -v -cover ./...

build:
	@echo "构建所有服务..."
	@mkdir -p $(OUTPUT_DIR)
	@echo "构建用户API服务..."
	go build -o $(OUTPUT_DIR)/user-api $(ROOT_DIR)/apps/user/api/user.go
	@echo "构建用户RPC服务..."
	go build -o $(OUTPUT_DIR)/user-rpc $(ROOT_DIR)/apps/user/rpc/rpc.go
	@echo "构建完成！二进制文件位于: $(OUTPUT_DIR)/"

# 开发工具安装
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/air-verse/air@latest

# 开发环境管理
dev:
	@echo "后台启动开发环境..."
	cd $(ROOT_DIR)/deploy/dev && docker compose up --build -d

dev-fg:
	@echo "前台启动开发环境..."
	cd $(ROOT_DIR)/deploy/dev && docker compose up --build

dev-stop:
	@echo "停止开发环境..."
	cd $(ROOT_DIR)/deploy/dev && docker compose down

dev-clean:
	@echo "清理开发环境（删除容器、网络、镜像和数据卷）..."
	cd $(ROOT_DIR)/deploy/dev && docker compose down -v --rmi all --remove-orphans

dev-logs:
	@echo "查看开发环境日志..."
	cd $(ROOT_DIR)/deploy/dev && docker compose logs -f

dev-restart:
	@echo "重启开发环境..."
	cd $(ROOT_DIR)/deploy/dev && docker compose restart

dev-status:
	@echo "查看服务状态..."
	cd $(ROOT_DIR)/deploy/dev && docker compose ps

env-start:
	@echo "创建共享网络..."
	docker network create miniblog-network || true
	@echo "启动基础环境服务（MySQL、Redis、Kafka）..."
	cd $(ROOT_DIR)/deploy/dev && docker compose -f docker-compose.env.yml up -d

env-stop:
	@echo "停止基础环境服务（MySQL、Redis、Kafka）..."
	cd $(ROOT_DIR)/deploy/dev && docker compose -f docker-compose.env.yml down

env-clean:
	@echo "清理基础环境和网络..."
	cd $(ROOT_DIR)/deploy/dev && docker compose -f docker-compose.env.yml down -v
	docker network rm miniblog-network || true


# 清理构建产物
clean:
	@echo "清理构建产物..."
	rm -rf $(OUTPUT_DIR)
	rm -rf $(ROOT_DIR)/tmp
	@echo "清理完成！"

add-copyright: # 添加版权头信息.
	@addlicense -v -f $(ROOT_DIR)/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,vendor,docs.deploy,logs,tmp,$(OUTPUT_DIR) --skip-files=\.pb\.go$