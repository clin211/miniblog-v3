#!/bin/bash

# 生产环境部署脚本
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi

    log_info "Docker 环境检查通过"
}

# 创建必要的目录
create_directories() {
    log_info "创建必要的目录..."

    mkdir -p logs/user-api
    mkdir -p logs/user-rpc
    mkdir -p nginx/conf.d/services

    log_info "目录创建完成"
}

# 复制配置文件
copy_configs() {
    log_info "检查配置文件..."

    # 检查 Nginx 配置文件（用于外部 nginx）
    if [ -f "./nginx/conf.d/miniblog.conf" ]; then
        log_info "Nginx 主配置文件已存在"
    else
        log_warn "请确保 nginx/conf.d/miniblog.conf 配置文件存在"
    fi

    # 检查服务配置文件
    if [ -d "./nginx/conf.d/services" ]; then
        service_count=$(ls ./nginx/conf.d/services/*.conf 2>/dev/null | wc -l)
        log_info "发现 $service_count 个服务配置文件"
    else
        log_warn "服务配置目录不存在"
    fi

    log_info "配置文件检查完成"
}

# 加载环境变量
load_env() {
    if [ -f ".env" ]; then
        log_info "加载环境变量文件..."
        export $(cat .env | grep -v '^#' | xargs)
        log_info "环境变量加载完成"
    else
        log_warn "未找到 .env 文件，使用默认配置"
    fi
}

# 部署服务
deploy_services() {
    log_info "开始部署应用服务..."

    # 检查基础设施服务是否运行
    check_infrastructure_services

    # 检查本地镜像是否存在
    check_local_images

    # 停止现有应用服务
    docker-compose down --remove-orphans

    # 启动应用服务（使用本地镜像）
    docker-compose up -d

    log_info "应用服务部署完成"
}

# 检查基础设施服务
check_infrastructure_services() {
    log_info "检查基础设施服务状态..."

    # 检查 MySQL
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-mysql"; then
        log_error "MySQL 服务未运行，请先启动基础设施服务"
        log_info "使用命令: ./infrastructure-manager.sh start"
        exit 1
    fi

    # 检查 Redis
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-redis"; then
        log_error "Redis 服务未运行，请先启动基础设施服务"
        log_info "使用命令: ./infrastructure-manager.sh start"
        exit 1
    fi

    # 检查 etcd
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-etcd"; then
        log_error "etcd 服务未运行，请先启动基础设施服务"
        log_info "使用命令: ./infrastructure-manager.sh start"
        exit 1
    fi

    # 检查 Zookeeper
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-zookeeper"; then
        log_error "Zookeeper 服务未运行，请先启动基础设施服务"
        log_info "使用命令: ./infrastructure-manager.sh start"
        exit 1
    fi

    # 检查 Kafka
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-kafka"; then
        log_error "Kafka 服务未运行，请先启动基础设施服务"
        log_info "使用命令: ./infrastructure-manager.sh start"
        exit 1
    fi

    log_info "基础设施服务检查通过"
}

# 检查本地镜像
check_local_images() {
    log_info "检查本地镜像..."

    # 检查 user-api 镜像
    if ! docker images | grep -q "miniblog-user-api.*release"; then
        log_error "本地镜像 miniblog-user-api:release 不存在"
        log_info "请确保镜像已通过 docker load 加载到本地"
        exit 1
    fi

    # 检查 user-rpc 镜像
    if ! docker images | grep -q "miniblog-user-rpc.*release"; then
        log_error "本地镜像 miniblog-user-rpc:release 不存在"
        log_info "请确保镜像已通过 docker load 加载到本地"
        exit 1
    fi

    log_info "本地镜像检查通过"
    docker images | grep miniblog
}

# 检查服务状态
check_services() {
    log_info "检查服务状态..."

    sleep 10

    # 检查容器状态
    docker-compose ps

    # 检查服务健康状态
    log_info "检查服务健康状态..."

    # 检查 user-rpc 服务
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_info "user-rpc 服务健康检查通过"
    else
        log_warn "user-rpc 服务健康检查失败"
    fi

    # 检查 user-api 服务
    if curl -f http://localhost:8888/health > /dev/null 2>&1; then
        log_info "user-api 服务健康检查通过"
    else
        log_warn "user-api 服务健康检查失败"
    fi

    # 检查外部 nginx 服务（如果配置了的话）
    if curl -f http://localhost/health > /dev/null 2>&1; then
        log_info "外部 nginx 服务健康检查通过"
    else
        log_warn "外部 nginx 服务健康检查失败（可能需要手动配置）"
    fi
}

# 显示服务信息
show_info() {
    log_info "服务部署信息："
    echo "=================================="
    echo "服务访问地址："
    echo "  - User API: http://localhost:8888"
    echo "  - User RPC: http://localhost:8080"
    echo "  - 外部 Nginx: 需要手动配置"
    echo ""
    echo "基础设施服务："
    echo "  - MySQL: localhost:3306"
    echo "  - Redis: localhost:6379"
    echo "  - etcd: localhost:2379"
    echo "  - Zookeeper: localhost:2181"
    echo "  - Kafka: localhost:29092"
    echo ""
    echo "日志目录："
    echo "  - User API: ./logs/user-api"
    echo "  - User RPC: ./logs/user-rpc"
    echo ""
    echo "管理脚本："
    echo "  - 基础设施管理: ./infrastructure-manager.sh"
    echo "  - Kafka 管理: ./kafka-manager.sh"
    echo "  - Nginx 管理: ./nginx-manager.sh"
    echo ""
    echo "下一步操作："
    echo "  1. 配置外部 nginx: ./nginx-manager.sh install"
    echo "  2. 创建 Kafka 主题: ./kafka-manager.sh create-default"
    echo "  3. 检查服务状态: ./infrastructure-manager.sh status"
    echo "  4. 查看应用日志: docker-compose logs -f user-api user-rpc"
    echo "=================================="
}

# 主函数
main() {
    log_info "开始生产环境部署..."

    check_docker
    create_directories
    copy_configs
    load_env
    deploy_services
    check_services
    show_info

    log_info "部署完成！"
}

# 执行主函数
main "$@"
