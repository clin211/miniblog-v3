#!/bin/bash

# 基础设施管理脚本
# 用于管理 MySQL、Redis 等持久化服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 配置变量
COMPOSE_FILE="docker-compose.env.yml"
PROJECT_NAME="miniblog-infrastructure"

# 启动基础设施服务
start_infrastructure() {
    log_info "启动基础设施服务..."

    # 检查 Docker 是否运行
    if ! docker info > /dev/null 2>&1; then
        log_error "Docker 未运行，请先启动 Docker"
        exit 1
    fi

    # 启动服务
    docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME up -d

    log_info "基础设施服务启动完成"
    log_info "等待服务就绪..."

    # 等待服务就绪
    wait_for_services

    log_info "所有基础设施服务已就绪"
}

# 停止基础设施服务
stop_infrastructure() {
    log_info "停止基础设施服务..."

    docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME down

    log_info "基础设施服务已停止"
}

# 重启基础设施服务
restart_infrastructure() {
    log_info "重启基础设施服务..."

    stop_infrastructure
    sleep 2
    start_infrastructure
}

# 查看基础设施服务状态
status_infrastructure() {
    log_info "基础设施服务状态："
    echo "=================================="

    docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME ps

    echo ""
    log_info "服务健康状态："

        # 检查 MySQL
    if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME ps mysql | grep -q "Up"; then
        log_info "MySQL: 运行中"
    else
        log_warn "MySQL: 未运行"
    fi

    # 检查 Redis
    if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME ps redis | grep -q "Up"; then
        log_info "Redis: 运行中"
    else
        log_warn "Redis: 未运行"
    fi

    # 检查 etcd
    if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME ps etcd | grep -q "Up"; then
        log_info "etcd: 运行中"
    else
        log_warn "etcd: 未运行"
    fi

    # 检查 Zookeeper
    if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME ps zookeeper | grep -q "Up"; then
        log_info "Zookeeper: 运行中"
    else
        log_warn "Zookeeper: 未运行"
    fi

    # 检查 Kafka
    if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME ps kafka | grep -q "Up"; then
        log_info "Kafka: 运行中"
    else
        log_warn "Kafka: 未运行"
    fi

    echo "=================================="
}

# 等待服务就绪
wait_for_services() {
    log_info "等待 MySQL 就绪..."
    local mysql_ready=false
    local attempts=0
    local max_attempts=30

    while [ "$mysql_ready" = false ] && [ $attempts -lt $max_attempts ]; do
        if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME exec -T mysql mysqladmin ping -h localhost --silent; then
            mysql_ready=true
            log_info "MySQL 已就绪"
        else
            attempts=$((attempts + 1))
            sleep 2
        fi
    done

    if [ "$mysql_ready" = false ]; then
        log_error "MySQL 启动超时"
        exit 1
    fi

    log_info "等待 Redis 就绪..."
    local redis_ready=false
    attempts=0

    while [ "$redis_ready" = false ] && [ $attempts -lt $max_attempts ]; do
        if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME exec -T redis redis-cli ping > /dev/null 2>&1; then
            redis_ready=true
            log_info "Redis 已就绪"
        else
            attempts=$((attempts + 1))
            sleep 2
        fi
    done

    if [ "$redis_ready" = false ]; then
        log_error "Redis 启动超时"
        exit 1
    fi

    log_info "等待 etcd 就绪..."
    local etcd_ready=false
    attempts=0

    while [ "$etcd_ready" = false ] && [ $attempts -lt $max_attempts ]; do
        if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME exec -T etcd etcdctl endpoint health > /dev/null 2>&1; then
            etcd_ready=true
            log_info "etcd 已就绪"
        else
            attempts=$((attempts + 1))
            sleep 2
        fi
    done

    if [ "$etcd_ready" = false ]; then
        log_error "etcd 启动超时"
        exit 1
    fi

    log_info "等待 Zookeeper 就绪..."
    local zookeeper_ready=false
    attempts=0

    while [ "$zookeeper_ready" = false ] && [ $attempts -lt $max_attempts ]; do
        if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME exec -T zookeeper zookeeper-shell localhost:2181 ruok | grep -q "imok"; then
            zookeeper_ready=true
            log_info "Zookeeper 已就绪"
        else
            attempts=$((attempts + 1))
            sleep 2
        fi
    done

    if [ "$zookeeper_ready" = false ]; then
        log_error "Zookeeper 启动超时"
        exit 1
    fi

    log_info "等待 Kafka 就绪..."
    local kafka_ready=false
    attempts=0

    while [ "$kafka_ready" = false ] && [ $attempts -lt $max_attempts ]; do
        if docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME exec -T kafka kafka-topics.sh --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
            kafka_ready=true
            log_info "Kafka 已就绪"
        else
            attempts=$((attempts + 1))
            sleep 5
        fi
    done

    if [ "$kafka_ready" = false ]; then
        log_error "Kafka 启动超时"
        exit 1
    fi
}

# 备份数据库
backup_database() {
    local backup_dir=${1:-"./backups"}
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="$backup_dir/mysql_backup_$timestamp.sql"

    log_info "备份 MySQL 数据库..."

    # 创建备份目录
    mkdir -p "$backup_dir"

    # 执行备份
    docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME exec -T mysql mysqldump \
        -u root -p${MYSQL_ROOT_PASSWORD:-root123456} \
        --all-databases > "$backup_file"

    if [ $? -eq 0 ]; then
        log_info "数据库备份完成: $backup_file"
    else
        log_error "数据库备份失败"
        exit 1
    fi
}

# 恢复数据库
restore_database() {
    local backup_file=$1

    if [ -z "$backup_file" ]; then
        log_error "请指定备份文件路径"
        log_info "用法: $0 restore <backup-file>"
        exit 1
    fi

    if [ ! -f "$backup_file" ]; then
        log_error "备份文件不存在: $backup_file"
        exit 1
    fi

    log_info "恢复 MySQL 数据库..."
    log_warn "这将覆盖现有数据，请确认是否继续 (y/N)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
        docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME exec -T mysql mysql \
            -u root -p${MYSQL_ROOT_PASSWORD:-root123456} < "$backup_file"

        if [ $? -eq 0 ]; then
            log_info "数据库恢复完成"
        else
            log_error "数据库恢复失败"
            exit 1
        fi
    else
        log_info "取消恢复操作"
    fi
}

# 查看日志
show_logs() {
    local service=${1:-""}

    if [ -z "$service" ]; then
        log_info "显示所有服务日志："
        docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME logs -f
    else
        log_info "显示 $service 服务日志："
        docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME logs -f "$service"
    fi
}

# 清理数据（危险操作）
clean_data() {
    log_error "这将删除所有数据，包括数据库和缓存"
    log_warn "请确认是否继续 (y/N)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
        log_info "停止服务..."
        docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME down

            log_info "删除数据卷..."
    docker volume rm ${PROJECT_NAME}_mysql_data ${PROJECT_NAME}_redis_data ${PROJECT_NAME}_etcd_data ${PROJECT_NAME}_zookeeper_data ${PROJECT_NAME}_zookeeper_logs ${PROJECT_NAME}_kafka_data 2>/dev/null || true

        log_info "数据清理完成"
    else
        log_info "取消清理操作"
    fi
}

# 显示帮助信息
show_help() {
    echo "基础设施管理脚本"
    echo ""
    echo "用法: $0 <命令> [参数]"
    echo ""
    echo "命令:"
    echo "  start                   启动基础设施服务"
    echo "  stop                    停止基础设施服务"
    echo "  restart                 重启基础设施服务"
    echo "  status                  查看服务状态"
    echo "  logs [service]          查看服务日志"
    echo "  backup [backup-dir]     备份数据库"
    echo "  restore <backup-file>   恢复数据库"
    echo "  clean                   清理所有数据（危险操作）"
    echo "  help                    显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 start                 # 启动基础设施"
    echo "  $0 status                # 查看状态"
    echo "  $0 backup ./backups      # 备份数据库"
    echo "  $0 restore backup.sql    # 恢复数据库"
    echo "  $0 logs mysql            # 查看 MySQL 日志"
}

# 主函数
main() {
    case "${1:-help}" in
        "start")
            start_infrastructure
            ;;
        "stop")
            stop_infrastructure
            ;;
        "restart")
            restart_infrastructure
            ;;
        "status")
            status_infrastructure
            ;;
        "logs")
            show_logs "$2"
            ;;
        "backup")
            backup_database "$2"
            ;;
        "restore")
            restore_database "$2"
            ;;
        "clean")
            clean_data
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
