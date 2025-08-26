#!/bin/bash

# 生产环境部署脚本
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
IMAGES_DIR="/home/project/miniblog-v3/images"
USER_API_IMAGE="miniblog-v3-user-api:release"
USER_RPC_IMAGE="miniblog-v3-user-rpc:release"

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

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
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

# 检查镜像文件是否存在
check_image_files() {
    log_info "检查镜像文件..."

    local user_api_tar="$IMAGES_DIR/miniblog-v3-user-api.tar"
    local user_rpc_tar="$IMAGES_DIR/miniblog-v3-user-rpc.tar"

    if [ ! -f "$user_api_tar" ]; then
        log_error "镜像文件不存在: $user_api_tar"
        log_info "请确保镜像文件已上传到 $IMAGES_DIR 目录"
        exit 1
    fi

    if [ ! -f "$user_rpc_tar" ]; then
        log_error "镜像文件不存在: $user_rpc_tar"
        log_info "请确保镜像文件已上传到 $IMAGES_DIR 目录"
        exit 1
    fi

    log_info "镜像文件检查通过"
    log_info "发现镜像文件："
    ls -lh "$IMAGES_DIR"/*.tar
}

# 加载镜像到 Docker
load_images() {
    log_info "开始加载镜像到 Docker..."

    local user_api_tar="$IMAGES_DIR/miniblog-v3-user-api.tar"
    local user_rpc_tar="$IMAGES_DIR/miniblog-v3-user-rpc.tar"

    # 检查镜像是否已存在
    if docker images | grep -q "$USER_API_IMAGE"; then
        log_warn "镜像 $USER_API_IMAGE 已存在，跳过加载"
    else
        log_info "加载 user-api 镜像..."
        if docker load -i "$user_api_tar"; then
            log_info "user-api 镜像加载成功"
        else
            log_error "user-api 镜像加载失败"
            exit 1
        fi
    fi

    if docker images | grep -q "$USER_RPC_IMAGE"; then
        log_warn "镜像 $USER_RPC_IMAGE 已存在，跳过加载"
    else
        log_info "加载 user-rpc 镜像..."
        if docker load -i "$user_rpc_tar"; then
            log_info "user-rpc 镜像加载成功"
        else
            log_error "user-rpc 镜像加载失败"
            exit 1
        fi
    fi

    log_info "镜像加载完成"
}

# 部署服务
deploy_services() {
    log_info "开始部署应用服务..."

    # 检查基础设施服务是否运行
    check_infrastructure_services

    # 检查并加载镜像
    check_image_files
    check_local_images

    # 停止现有应用服务
    log_info "停止现有应用服务..."
    docker-compose down --remove-orphans

    # 启动应用服务（使用本地镜像）
    log_info "启动应用服务..."
    docker-compose up -d

    log_info "应用服务部署完成"
}

# 检查基础设施服务
check_infrastructure_services() {
    log_info "检查基础设施服务状态..."

    local missing_services=()

    # 检查 MySQL
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-v3-mysql-1"; then
        missing_services+=("MySQL")
    fi

    # 检查 Redis
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-v3-redis-1"; then
        missing_services+=("Redis")
    fi

    # 检查 etcd
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-v3-etcd-1"; then
        missing_services+=("etcd")
    fi

    # 检查 Zookeeper
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-v3-zookeeper-1"; then
        missing_services+=("Zookeeper")
    fi

    # 检查 Kafka
    if ! docker ps --format "table {{.Names}}" | grep -q "miniblog-v3-kafka-1"; then
        missing_services+=("Kafka")
    fi

    if [ ${#missing_services[@]} -gt 0 ]; then
        log_error "以下基础设施服务未运行: ${missing_services[*]}"
        log_info "请先启动基础设施服务："
        log_info "  ./infrastructure-manager.sh start"
        log_info "或者使用 docker-compose 启动："
        log_info "  docker-compose -f docker-compose.env.yml up -d"
        exit 1
    fi

    log_info "基础设施服务检查通过"
}

# 检查本地镜像
check_local_images() {
    log_info "检查本地镜像..."

    # 检查 user-api 镜像
    if ! docker images | grep -q "$USER_API_IMAGE"; then
        log_error "本地镜像 $USER_API_IMAGE 不存在"
        log_info "尝试自动加载镜像..."
        load_images
    fi

    # 检查 user-rpc 镜像
    if ! docker images | grep -q "$USER_RPC_IMAGE"; then
        log_error "本地镜像 $USER_RPC_IMAGE 不存在"
        log_info "尝试自动加载镜像..."
        load_images
    fi

    log_info "本地镜像检查通过"
    log_info "发现的 miniblog-v3 镜像："
    docker images | grep miniblog-v3 || echo "没有找到 miniblog-v3 镜像"
}

# 检查服务状态
check_services() {
    log_info "检查服务状态..."

    # 等待服务启动
    log_info "等待服务启动..."
    sleep 15

    # 检查容器状态
    log_info "容器状态："
    docker-compose ps

    # 检查服务健康状态
    log_info "检查服务健康状态..."

    # 检查 user-rpc 服务
    local rpc_health_ok=false
    for i in {1..5}; do
        if curl -f http://localhost:8080/health > /dev/null 2>&1; then
            log_info "user-rpc 服务健康检查通过"
            rpc_health_ok=true
            break
        else
            log_warn "user-rpc 服务健康检查失败，重试 $i/5"
            sleep 3
        fi
    done

    if [ "$rpc_health_ok" = false ]; then
        log_error "user-rpc 服务健康检查最终失败"
        log_info "查看 user-rpc 日志："
        docker-compose logs --tail 20 user-rpc
    fi

    # 检查 user-api 服务
    local api_health_ok=false
    for i in {1..5}; do
        if curl -f http://localhost:8888/health > /dev/null 2>&1; then
            log_info "user-api 服务健康检查通过"
            api_health_ok=true
            break
        else
            log_warn "user-api 服务健康检查失败，重试 $i/5"
            sleep 3
        fi
    done

    if [ "$api_health_ok" = false ]; then
        log_error "user-api 服务健康检查最终失败"
        log_info "查看 user-api 日志："
        docker-compose logs --tail 20 user-api
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
    echo "  - MySQL: localhost:3306 (容器: miniblog-v3-mysql-1)"
    echo "  - Redis: localhost:6379 (容器: miniblog-v3-redis-1)"
    echo "  - etcd: localhost:2379 (容器: miniblog-v3-etcd-1)"
    echo "  - Zookeeper: localhost:2181 (容器: miniblog-v3-zookeeper-1)"
    echo "  - Kafka: localhost:29092 (容器: miniblog-v3-kafka-1)"
    echo ""
    echo "应用服务："
    echo "  - User API: miniblog-user-api (端口: 8888)"
    echo "  - User RPC: miniblog-user-rpc (端口: 8080)"
    echo ""
    echo "镜像信息："
    echo "  - 镜像目录: $IMAGES_DIR"
    echo "  - User API 镜像: $USER_API_IMAGE"
    echo "  - User RPC 镜像: $USER_RPC_IMAGE"
    echo ""
    echo "日志目录："
    echo "  - User API: ./logs/user-api"
    echo "  - User RPC: ./logs/user-rpc"
    echo ""
    echo "管理脚本："
    echo "  - 基础设施管理: ./infrastructure-manager.sh"
    echo "  - Kafka 管理: ./kafka-manager.sh"
    echo "  - Nginx 管理: ./nginx-manager.sh"
    echo "  - 健康检查: ./health-check.sh"
    echo ""
    echo "下一步操作："
    echo "  1. 配置外部 nginx: ./nginx-manager.sh install"
    echo "  2. 创建 Kafka 主题: ./kafka-manager.sh create-default"
    echo "  3. 检查服务状态: ./infrastructure-manager.sh status"
    echo "  4. 健康检查: ./infrastructure-manager.sh health all"
    echo "  5. 查看应用日志: docker-compose logs -f user-api user-rpc"
    echo "  6. 重新加载镜像: 删除镜像后重新运行此脚本"
    echo "=================================="
}

# 清理函数
cleanup() {
    log_info "清理临时资源..."
    # 可以在这里添加清理逻辑
}

# 主函数
main() {
    log_info "开始生产环境部署..."

    # 设置错误处理
    trap cleanup EXIT

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
