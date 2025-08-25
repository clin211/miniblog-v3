#!/bin/bash

# Kafka 管理脚本
# 用于管理 Kafka 主题和配置

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
KAFKA_CONTAINER="miniblog-kafka"
BOOTSTRAP_SERVER="localhost:29092"

# 检查 Kafka 容器是否运行
check_kafka_container() {
    if ! docker ps --format "table {{.Names}}" | grep -q "$KAFKA_CONTAINER"; then
        log_error "Kafka 容器未运行，请先启动基础设施服务"
        log_info "使用命令: ./infrastructure-manager.sh start"
        exit 1
    fi
}

# 列出所有主题
list_topics() {
    log_info "列出所有 Kafka 主题："
    echo "=================================="

          docker exec $KAFKA_CONTAINER kafka-topics.sh \
          --bootstrap-server localhost:9092 \
          --list

    echo "=================================="
}

# 创建主题
create_topic() {
    local topic_name=$1
    local partitions=${2:-1}
    local replication_factor=${3:-1}

    if [ -z "$topic_name" ]; then
        log_error "请指定主题名称"
        log_info "用法: $0 create <topic-name> [partitions] [replication-factor]"
        exit 1
    fi

    log_info "创建主题: $topic_name (分区: $partitions, 副本因子: $replication_factor)"

          docker exec $KAFKA_CONTAINER kafka-topics.sh \
          --bootstrap-server localhost:9092 \
          --create \
          --topic "$topic_name" \
          --partitions "$partitions" \
          --replication-factor "$replication_factor"

    log_info "主题创建完成"
}

# 删除主题
delete_topic() {
    local topic_name=$1

    if [ -z "$topic_name" ]; then
        log_error "请指定主题名称"
        log_info "用法: $0 delete <topic-name>"
        exit 1
    fi

    log_warn "删除主题: $topic_name"
    log_warn "这将永久删除主题及其所有数据，请确认是否继续 (y/N)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
                  docker exec $KAFKA_CONTAINER kafka-topics.sh \
              --bootstrap-server localhost:9092 \
              --delete \
              --topic "$topic_name"

        log_info "主题删除完成"
    else
        log_info "取消删除操作"
    fi
}

# 查看主题详情
describe_topic() {
    local topic_name=$1

    if [ -z "$topic_name" ]; then
        log_error "请指定主题名称"
        log_info "用法: $0 describe <topic-name>"
        exit 1
    fi

    log_info "主题详情: $topic_name"
    echo "=================================="

          docker exec $KAFKA_CONTAINER kafka-topics.sh \
          --bootstrap-server localhost:9092 \
          --describe \
          --topic "$topic_name"

    echo "=================================="
}

# 查看消费者组
list_consumer_groups() {
    log_info "列出所有消费者组："
    echo "=================================="

          docker exec $KAFKA_CONTAINER kafka-consumer-groups.sh \
          --bootstrap-server localhost:9092 \
          --list

    echo "=================================="
}

# 查看消费者组详情
describe_consumer_group() {
    local group_name=$1

    if [ -z "$group_name" ]; then
        log_error "请指定消费者组名称"
        log_info "用法: $0 describe-group <group-name>"
        exit 1
    fi

    log_info "消费者组详情: $group_name"
    echo "=================================="

          docker exec $KAFKA_CONTAINER kafka-consumer-groups.sh \
          --bootstrap-server localhost:9092 \
          --describe \
          --group "$group_name"

    echo "=================================="
}

# 重置消费者组偏移量
reset_consumer_group() {
    local group_name=$1
    local topic_name=$2
    local offset=${3:-"earliest"}

    if [ -z "$group_name" ] || [ -z "$topic_name" ]; then
        log_error "请指定消费者组名称和主题名称"
        log_info "用法: $0 reset-group <group-name> <topic-name> [earliest|latest]"
        exit 1
    fi

    log_warn "重置消费者组偏移量: $group_name -> $topic_name ($offset)"
    log_warn "这将重置消费者组的偏移量，请确认是否继续 (y/N)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
                  docker exec $KAFKA_CONTAINER kafka-consumer-groups.sh \
              --bootstrap-server localhost:9092 \
              --group "$group_name" \
              --topic "$topic_name" \
              --reset-offsets \
              --to-"$offset" \
              --execute

        log_info "偏移量重置完成"
    else
        log_info "取消重置操作"
    fi
}

# 发送测试消息
send_test_message() {
    local topic_name=$1
    local message=${2:-"Hello Kafka!"}

    if [ -z "$topic_name" ]; then
        log_error "请指定主题名称"
        log_info "用法: $0 send <topic-name> [message]"
        exit 1
    fi

    log_info "发送测试消息到主题: $topic_name"

          echo "$message" | docker exec -i $KAFKA_CONTAINER kafka-console-producer.sh \
          --bootstrap-server localhost:9092 \
          --topic "$topic_name"

    log_info "消息发送完成"
}

# 消费消息
consume_messages() {
    local topic_name=$1
    local group_name=${2:-"test-consumer-group"}
    local offset=${3:-"earliest"}

    if [ -z "$topic_name" ]; then
        log_error "请指定主题名称"
        log_info "用法: $0 consume <topic-name> [group-name] [earliest|latest]"
        exit 1
    fi

    log_info "消费主题消息: $topic_name (消费者组: $group_name, 偏移量: $offset)"
    log_info "按 Ctrl+C 停止消费"

          docker exec $KAFKA_CONTAINER kafka-console-consumer.sh \
          --bootstrap-server localhost:9092 \
          --topic "$topic_name" \
          --group "$group_name" \
          --from-beginning
}

# 查看集群信息
cluster_info() {
    log_info "Kafka 集群信息："
    echo "=================================="

          log_info "Broker 信息："
      docker exec $KAFKA_CONTAINER kafka-broker-api-versions.sh \
          --bootstrap-server localhost:9092

    echo ""
          log_info "配置信息："
      docker exec $KAFKA_CONTAINER kafka-configs.sh \
          --bootstrap-server localhost:9092 \
          --entity-type brokers \
          --entity-name 1 \
          --describe

    echo "=================================="
}

# 创建默认主题
create_default_topics() {
    log_info "创建默认主题..."

    # 用户相关主题
    create_topic "user-events" 3 1
    create_topic "user-notifications" 3 1

    # 博客相关主题
    create_topic "blog-events" 3 1
    create_topic "blog-notifications" 3 1

    # 系统主题
    create_topic "system-events" 1 1
    create_topic "audit-logs" 3 1

    log_info "默认主题创建完成"
}

# 显示帮助信息
show_help() {
    echo "Kafka 管理脚本"
    echo ""
    echo "用法: $0 <命令> [参数]"
    echo ""
    echo "命令:"
    echo "  list                          列出所有主题"
    echo "  create <topic> [partitions] [replicas]  创建主题"
    echo "  delete <topic>                删除主题"
    echo "  describe <topic>              查看主题详情"
    echo "  list-groups                   列出所有消费者组"
    echo "  describe-group <group>        查看消费者组详情"
    echo "  reset-group <group> <topic> [offset]  重置消费者组偏移量"
    echo "  send <topic> [message]        发送测试消息"
    echo "  consume <topic> [group] [offset]  消费消息"
    echo "  cluster                       查看集群信息"
    echo "  create-default                创建默认主题"
    echo "  help                          显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 list                        # 列出所有主题"
    echo "  $0 create user-events 3 1      # 创建主题"
    echo "  $0 describe user-events        # 查看主题详情"
    echo "  $0 send user-events 'Hello'    # 发送消息"
    echo "  $0 consume user-events         # 消费消息"
}

# 主函数
main() {
    check_kafka_container

    case "${1:-help}" in
        "list")
            list_topics
            ;;
        "create")
            create_topic "$2" "$3" "$4"
            ;;
        "delete")
            delete_topic "$2"
            ;;
        "describe")
            describe_topic "$2"
            ;;
        "list-groups")
            list_consumer_groups
            ;;
        "describe-group")
            describe_consumer_group "$2"
            ;;
        "reset-group")
            reset_consumer_group "$2" "$3" "$4"
            ;;
        "send")
            send_test_message "$2" "$3"
            ;;
        "consume")
            consume_messages "$2" "$3" "$4"
            ;;
        "cluster")
            cluster_info
            ;;
        "create-default")
            create_default_topics
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
