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

    # 创建网络
    log_info "创建网络 miniblog-v3-network..."
    docker network create miniblog-v3-network 2>/dev/null || log_info "网络已存在"

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

    # 删除网络
    log_info "删除网络 miniblog-v3-network..."
    docker network rm miniblog-v3-network 2>/dev/null || log_info "网络不存在或已被删除"

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

    # 检查网络状态
    if docker network ls | grep -q "miniblog-v3-network"; then
        log_info "网络 miniblog-v3-network: 已创建"
    else
        log_warn "网络 miniblog-v3-network: 未创建"
    fi

    echo ""
    docker-compose -f $COMPOSE_FILE -p $PROJECT_NAME ps

    echo ""
    log_info "服务健康状态："

        # 检查 MySQL
    if docker ps | grep -q "miniblog-v3-mysql-1"; then
        log_info "MySQL: 运行中"
    else
        log_warn "MySQL: 未运行"
    fi

    # 检查 Redis
    if docker ps | grep -q "miniblog-v3-redis-1"; then
        log_info "Redis: 运行中"
    else
        log_warn "Redis: 未运行"
    fi

    # 检查 etcd
    if docker ps | grep -q "miniblog-v3-etcd-1"; then
        log_info "etcd: 运行中"
    else
        log_warn "etcd: 未运行"
    fi

    # 检查 Zookeeper
    if docker ps | grep -q "miniblog-v3-zookeeper-1"; then
        log_info "Zookeeper: 运行中"
    else
        log_warn "Zookeeper: 未运行"
    fi

    # 检查 Kafka
    if docker ps | grep -q "miniblog-v3-kafka-1"; then
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
        if docker exec miniblog-v3-mysql-1 mysqladmin ping -h localhost --silent; then
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
        if docker exec miniblog-v3-redis-1 redis-cli -a "redis123" ping > /dev/null 2>&1; then
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
        if docker exec miniblog-v3-etcd-1 etcdctl endpoint health --endpoints=http://localhost:2379 > /dev/null 2>&1; then
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
        if docker exec miniblog-v3-zookeeper-1 sh -c "echo ruok | nc localhost 2181 | grep -q imok"; then
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
        if docker exec miniblog-v3-kafka-1 kafka-topics.sh --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
            kafka_ready=true
            log_info "Kafka 已就绪"
        else
            attempts=$((attempts + 1))
            sleep 10
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
    docker exec miniblog-v3-mysql-1 mysqldump \
        -u root -proot123456 \
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
        docker exec -i miniblog-v3-mysql-1 mysql \
            -u root -proot123456 < "$backup_file"

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

# 测试服务连接
test_service() {
    local service=${1:-""}

    if [ -z "$service" ]; then
        log_error "请指定要测试的服务"
        log_info "可用服务: mysql, redis, etcd, zookeeper, kafka"
        exit 1
    fi

    case "$service" in
        "mysql")
            log_info "测试 MySQL 连接..."
            if docker exec miniblog-v3-mysql-1 mysqladmin ping -h localhost --silent; then
                log_info "MySQL 连接正常"
            else
                log_error "MySQL 连接失败"
            fi
            ;;
        "redis")
            log_info "测试 Redis 连接..."
            if docker exec miniblog-v3-redis-1 redis-cli -a "redis123" ping > /dev/null 2>&1; then
                log_info "Redis 连接正常"
            else
                log_error "Redis 连接失败"
            fi
            ;;
        "etcd")
            log_info "测试 etcd 连接..."
            if docker exec miniblog-v3-etcd-1 etcdctl endpoint health --endpoints=http://localhost:2379 > /dev/null 2>&1; then
                log_info "etcd 连接正常"
            else
                log_error "etcd 连接失败"
            fi
            ;;
        "zookeeper")
            log_info "测试 Zookeeper 连接..."
            if docker exec miniblog-v3-zookeeper-1 sh -c "echo ruok | nc localhost 2181 | grep -q imok"; then
                log_info "Zookeeper 连接正常"
            else
                log_error "Zookeeper 连接失败"
            fi
            ;;
        "kafka")
            log_info "测试 Kafka 连接..."
            if docker exec miniblog-v3-kafka-1 kafka-topics.sh --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
                log_info "Kafka 连接正常"
            else
                log_error "Kafka 连接失败"
            fi
            ;;
        *)
            log_error "未知服务: $service"
            log_info "可用服务: mysql, redis, etcd, zookeeper, kafka"
            exit 1
            ;;
    esac
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
        docker volume rm miniblog-infrastructure_mysql_data miniblog-infrastructure_redis_data miniblog-infrastructure_etcd_data miniblog-infrastructure_zookeeper_data miniblog-infrastructure_zookeeper_logs miniblog-infrastructure_kafka_data 2>/dev/null || true

        log_info "删除网络..."
        docker network rm miniblog-v3-network 2>/dev/null || true

        log_info "数据清理完成"
    else
        log_info "取消清理操作"
    fi
}

# 网络管理
manage_network() {
    local action=${1:-"status"}

    case "$action" in
        "create")
            log_info "创建网络 miniblog-v3-network..."
            docker network create miniblog-v3-network
            log_info "网络创建完成"
            ;;
        "delete")
            log_info "删除网络 miniblog-v3-network..."
            docker network rm miniblog-v3-network
            log_info "网络删除完成"
            ;;
        "status")
            if docker network ls | grep -q "miniblog-v3-network"; then
                log_info "网络 miniblog-v3-network 已存在"
                docker network inspect miniblog-v3-network --format='{{.Name}}: {{.Driver}}'
            else
                log_warn "网络 miniblog-v3-network 不存在"
            fi
            ;;
        *)
            log_error "未知网络操作: $action"
            log_info "可用操作: create, delete, status"
            exit 1
            ;;
    esac
}

# 健康检查
health_check() {
    local service=${1:-"all"}

    # 检查健康检查脚本是否存在
    if [ ! -f "./health-check.sh" ]; then
        log_error "健康检查脚本不存在: ./health-check.sh"
        exit 1
    fi

    # 确保脚本有执行权限
    chmod +x ./health-check.sh

    case "$service" in
        "all")
            log_info "执行所有服务健康检查..."
            ./health-check.sh check all
            ;;
        "mysql"|"redis"|"etcd"|"zookeeper"|"kafka")
            log_info "执行 $service 服务健康检查..."
            ./health-check.sh check "$service"
            ;;
        "quick")
            log_info "执行快速健康检查..."
            ./health-check.sh quick all
            ;;
        *)
            log_error "未知服务: $service"
            log_info "可用服务: mysql, redis, etcd, zookeeper, kafka, all, quick"
            exit 1
            ;;
    esac
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
    echo "  test <service>          测试服务连接"
    echo "  logs [service]          查看服务日志"
    echo "  backup [backup-dir]     备份数据库"
    echo "  restore <backup-file>   恢复数据库"
    echo "  clean                   清理所有数据（危险操作）"
    echo "  network <action>        网络管理 (create/delete/status)"
    echo "  health [service]        健康检查 (mysql/redis/etcd/zookeeper/kafka/all/quick)"
    echo "  help                    显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 start                 # 启动基础设施"
    echo "  $0 status                # 查看状态"
    echo "  $0 test zookeeper        # 测试 Zookeeper 连接"
    echo "  $0 backup ./backups      # 备份数据库"
    echo "  $0 restore backup.sql    # 恢复数据库"
    echo "  $0 logs mysql            # 查看 MySQL 日志"
    echo "  $0 network status        # 查看网络状态"
    echo "  $0 network create        # 创建网络"
    echo "  $0 network delete        # 删除网络"
    echo "  $0 health all            # 检查所有服务健康状态"
    echo "  $0 health mysql          # 检查 MySQL 健康状态"
    echo "  $0 health quick          # 快速健康检查（静默模式）"
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
        "test")
            test_service "$2"
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
        "network")
            manage_network "$2"
            ;;
        "health")
            health_check "$2"
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