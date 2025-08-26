#!/bin/bash

# 服务健康检查脚本
# 用于检查所有基础设施服务的健康状态

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 配置变量
MYSQL_ROOT_PASSWORD="root123456"
MYSQL_USER="miniblog"
MYSQL_PASSWORD="miniblog123"
REDIS_PASSWORD="redis123"

# 检查 MySQL 服务
check_mysql() {
    log_info "检查 MySQL 服务..."

    # 检查容器是否运行
    if ! docker ps | grep -q "miniblog-v3-mysql-1"; then
        log_error "MySQL 容器未运行"
        return 1
    fi

    # 检查 MySQL 连接
    if docker exec miniblog-v3-mysql-1 mysqladmin ping -h localhost --silent; then
        log_success "MySQL 服务正常"

        # 检查数据库是否存在
        if docker exec miniblog-v3-mysql-1 mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "USE miniblog;" 2>/dev/null; then
            log_success "MySQL 数据库 miniblog 存在"
        else
            log_warning "MySQL 数据库 miniblog 不存在"
        fi

        return 0
    else
        log_error "MySQL 服务异常"
        return 1
    fi
}

# 检查 Redis 服务
check_redis() {
    log_info "检查 Redis 服务..."

    # 检查容器是否运行
    if ! docker ps | grep -q "miniblog-v3-redis-1"; then
        log_error "Redis 容器未运行"
        return 1
    fi

    # 检查 Redis 连接（使用密码）
    if docker exec miniblog-v3-redis-1 redis-cli -a "${REDIS_PASSWORD}" ping 2>/dev/null | grep -q "PONG"; then
        log_success "Redis 服务正常"
        return 0
    else
        log_error "Redis 服务异常"
        return 1
    fi
}

# 检查 etcd 服务
check_etcd() {
    log_info "检查 etcd 服务..."

    # 检查容器是否运行
    if ! docker ps | grep -q "miniblog-v3-etcd-1"; then
        log_error "etcd 容器未运行"
        return 1
    fi

    # 检查 etcd 健康状态
    if docker exec miniblog-v3-etcd-1 etcdctl endpoint health --endpoints=http://localhost:2379 2>/dev/null | grep -q "healthy"; then
        log_success "etcd 服务正常"
        return 0
    else
        log_error "etcd 服务异常"
        return 1
    fi
}

# 检查 Zookeeper 服务
check_zookeeper() {
    log_info "检查 Zookeeper 服务..."

    # 检查容器是否运行
    if ! docker ps | grep -q "miniblog-v3-zookeeper-1"; then
        log_error "Zookeeper 容器未运行"
        return 1
    fi

    # 检查 Zookeeper 状态（使用四字命令）
    if docker exec miniblog-v3-zookeeper-1 sh -c "echo ruok | nc localhost 2181 | grep -q imok" 2>/dev/null; then
        log_success "Zookeeper 服务正常"
        return 0
    else
        log_error "Zookeeper 服务异常"
        return 1
    fi
}

# 检查 Kafka 服务
check_kafka() {
    log_info "检查 Kafka 服务..."

    # 检查容器是否运行
    if ! docker ps | grep -q "miniblog-v3-kafka-1"; then
        log_error "Kafka 容器未运行"
        return 1
    fi

    # 检查 Kafka 连接
    if docker exec miniblog-v3-kafka-1 kafka-topics.sh --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
        log_success "Kafka 服务正常"
        return 0
    else
        log_error "Kafka 服务异常"
        return 1
    fi
}

# 检查所有服务
check_all_services() {
    local failed_services=()

    echo "=================================="
    log_info "开始检查所有基础设施服务..."
    echo "=================================="

    # 检查 MySQL
    if check_mysql; then
        log_success "✅ MySQL 检查通过"
    else
        log_error "❌ MySQL 检查失败"
        failed_services+=("MySQL")
    fi

    # 检查 Redis
    if check_redis; then
        log_success "✅ Redis 检查通过"
    else
        log_error "❌ Redis 检查失败"
        failed_services+=("Redis")
    fi

    # 检查 etcd
    if check_etcd; then
        log_success "✅ etcd 检查通过"
    else
        log_error "❌ etcd 检查失败"
        failed_services+=("etcd")
    fi

    # 检查 Zookeeper
    if check_zookeeper; then
        log_success "✅ Zookeeper 检查通过"
    else
        log_error "❌ Zookeeper 检查失败"
        failed_services+=("Zookeeper")
    fi

    # 检查 Kafka
    if check_kafka; then
        log_success "✅ Kafka 检查通过"
    else
        log_error "❌ Kafka 检查失败"
        failed_services+=("Kafka")
    fi

    echo "=================================="

    # 输出检查结果
    if [ ${#failed_services[@]} -eq 0 ]; then
        log_success "🎉 所有基础设施服务运行正常！"
        return 0
    else
        log_error "❌ 以下服务存在问题："
        for service in "${failed_services[@]}"; do
            echo "  - $service"
        done
        log_info "请执行以下命令启动基础设施服务："
        echo "  ./infrastructure-manager.sh start"
        return 1
    fi
}

# 快速检查（用于 CI/CD）
quick_check() {
    local service=${1:-""}

    case "$service" in
        "mysql")
            check_mysql > /dev/null 2>&1
            return $?
            ;;
        "redis")
            check_redis > /dev/null 2>&1
            return $?
            ;;
        "etcd")
            check_etcd > /dev/null 2>&1
            return $?
            ;;
        "zookeeper")
            check_zookeeper > /dev/null 2>&1
            return $?
            ;;
        "kafka")
            check_kafka > /dev/null 2>&1
            return $?
            ;;
        "all")
            check_all_services > /dev/null 2>&1
            return $?
            ;;
        *)
            log_error "未知服务: $service"
            log_info "可用服务: mysql, redis, etcd, zookeeper, kafka, all"
            return 1
            ;;
    esac
}

# 显示帮助信息
show_help() {
    echo "服务健康检查脚本"
    echo ""
    echo "用法: $0 [命令] [服务名]"
    echo ""
    echo "命令:"
    echo "  check [service]    检查指定服务或所有服务"
    echo "  quick [service]    快速检查（静默模式，用于 CI/CD）"
    echo "  help               显示此帮助信息"
    echo ""
    echo "服务名:"
    echo "  mysql              检查 MySQL 服务"
    echo "  redis              检查 Redis 服务"
    echo "  etcd               检查 etcd 服务"
    echo "  zookeeper          检查 Zookeeper 服务"
    echo "  kafka              检查 Kafka 服务"
    echo "  all                检查所有服务"
    echo ""
    echo "示例:"
    echo "  $0 check all        # 检查所有服务"
    echo "  $0 check mysql      # 检查 MySQL 服务"
    echo "  $0 quick all        # 快速检查所有服务（静默模式）"
}

# 主函数
main() {
    case "${1:-help}" in
        "check")
            if [ -z "$2" ]; then
                check_all_services
            else
                case "$2" in
                    "mysql")
                        check_mysql
                        ;;
                    "redis")
                        check_redis
                        ;;
                    "etcd")
                        check_etcd
                        ;;
                    "zookeeper")
                        check_zookeeper
                        ;;
                    "kafka")
                        check_kafka
                        ;;
                    "all")
                        check_all_services
                        ;;
                    *)
                        log_error "未知服务: $2"
                        show_help
                        exit 1
                        ;;
                esac
            fi
            ;;
        "quick")
            quick_check "$2"
            exit $?
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
