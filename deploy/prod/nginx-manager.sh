#!/bin/bash

# Nginx 配置管理脚本
# 用于将 miniblog 配置添加到现有的 nginx 容器中

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

# 检查 nginx 容器是否存在
check_nginx_container() {
    log_info "检查 nginx 容器..."

    if ! docker ps --format "table {{.Names}}" | grep -q "nginx"; then
        log_error "未找到运行中的 nginx 容器"
        log_info "请确保 nginx 容器正在运行："
        log_info "docker ps | grep nginx"
        exit 1
    fi

    NGINX_CONTAINER=$(docker ps --format "table {{.Names}}" | grep "nginx" | head -1)
    log_info "找到 nginx 容器: $NGINX_CONTAINER"
}

# 备份现有配置
backup_nginx_config() {
    log_info "备份现有 nginx 配置..."

    BACKUP_DIR="/tmp/nginx_backup_$(date +%Y%m%d_%H%M%S)"
    docker exec $NGINX_CONTAINER mkdir -p $BACKUP_DIR

    # 备份现有配置文件
    docker exec $NGINX_CONTAINER cp -r /etc/nginx/conf.d $BACKUP_DIR/ || true
    docker exec $NGINX_CONTAINER cp /etc/nginx/nginx.conf $BACKUP_DIR/ || true

    log_info "配置已备份到容器内的 $BACKUP_DIR"
}

# 复制配置文件到 nginx 容器
copy_config_to_nginx() {
    log_info "复制 miniblog 配置到 nginx 容器..."

    # 复制主配置文件
    docker cp nginx/conf.d/miniblog.conf $NGINX_CONTAINER:/etc/nginx/conf.d/

    # 设置主配置文件权限
    docker exec $NGINX_CONTAINER chown root:root /etc/nginx/conf.d/miniblog.conf
    docker exec $NGINX_CONTAINER chmod 644 /etc/nginx/conf.d/miniblog.conf

    # 创建服务配置目录
    docker exec $NGINX_CONTAINER mkdir -p /etc/nginx/conf.d/services

    # 复制所有服务配置文件
    if [ -d "nginx/conf.d/services" ]; then
        for config_file in nginx/conf.d/services/*.conf; do
            if [ -f "$config_file" ]; then
                local service_name=$(basename "$config_file")
                docker cp "$config_file" "$NGINX_CONTAINER:/etc/nginx/conf.d/services/"
                docker exec $NGINX_CONTAINER chown root:root "/etc/nginx/conf.d/services/$service_name"
                docker exec $NGINX_CONTAINER chmod 644 "/etc/nginx/conf.d/services/$service_name"
                log_info "已复制服务配置: $service_name"
            fi
        done
    fi

    log_info "所有配置文件已复制到 nginx 容器"
}

# 测试 nginx 配置
test_nginx_config() {
    log_info "测试 nginx 配置..."

    if docker exec $NGINX_CONTAINER nginx -t; then
        log_info "nginx 配置测试通过"
    else
        log_error "nginx 配置测试失败"
        log_info "请检查配置文件并修复错误"
        exit 1
    fi
}

# 重新加载 nginx 配置
reload_nginx() {
    log_info "重新加载 nginx 配置..."

    # 发送 SIGHUP 信号重新加载配置
    docker exec $NGINX_CONTAINER nginx -s reload

    if [ $? -eq 0 ]; then
        log_info "nginx 配置重新加载成功"
    else
        log_error "nginx 配置重新加载失败"
        exit 1
    fi
}

# 验证配置是否生效
verify_config() {
    log_info "验证配置是否生效..."

    sleep 2

    # 检查 nginx 进程
    if docker exec $NGINX_CONTAINER pgrep nginx > /dev/null; then
        log_info "nginx 进程运行正常"
    else
        log_error "nginx 进程异常"
        exit 1
    fi

    # 检查配置文件是否加载
    if docker exec $NGINX_CONTAINER nginx -T | grep -q "miniblog.conf"; then
        log_info "miniblog 配置已加载"
    else
        log_warn "miniblog 配置可能未正确加载"
    fi
}

# 显示配置信息
show_config_info() {
    log_info "配置信息："
    echo "=================================="
    echo "Nginx 容器: $NGINX_CONTAINER"
    echo "配置文件: /etc/nginx/conf.d/miniblog.conf"
    echo ""
    echo "访问地址："
    echo "  - HTTPS: https://your-domain.com"
    echo "  - API: https://your-domain.com/api/"
    echo "  - 健康检查: https://your-domain.com/health"
    echo ""
    echo "注意事项："
    echo "  1. 请将配置文件中的 'your-domain.com' 替换为实际域名"
    echo "  2. 确保 SSL 证书路径正确"
    echo "  3. 确保防火墙允许 80 和 443 端口"
    echo "=================================="
}

# 恢复配置
restore_config() {
    log_info "恢复 nginx 配置..."

    read -p "请输入备份目录路径: " BACKUP_PATH

    if [ -z "$BACKUP_PATH" ]; then
        log_error "备份路径不能为空"
        exit 1
    fi

    # 恢复配置文件
    docker exec $NGINX_CONTAINER cp $BACKUP_PATH/nginx.conf /etc/nginx/ || true
    docker exec $NGINX_CONTAINER cp -r $BACKUP_PATH/conf.d/* /etc/nginx/conf.d/ || true

    # 重新加载配置
    reload_nginx

    log_info "配置已恢复"
}

# 显示帮助信息
show_help() {
    echo "Nginx 配置管理脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  install    安装 miniblog 配置到 nginx"
    echo "  restore    恢复 nginx 配置"
    echo "  test       测试 nginx 配置"
    echo "  reload     重新加载 nginx 配置"
    echo "  help       显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 install    # 安装配置"
    echo "  $0 restore    # 恢复配置"
}

# 主函数
main() {
    case "${1:-install}" in
        "install")
            log_info "开始安装 miniblog 配置到 nginx..."
            check_nginx_container
            backup_nginx_config
            copy_config_to_nginx
            test_nginx_config
            reload_nginx
            verify_config
            show_config_info
            log_info "配置安装完成！"
            ;;
        "restore")
            log_info "开始恢复 nginx 配置..."
            check_nginx_container
            restore_config
            ;;
        "test")
            log_info "测试 nginx 配置..."
            check_nginx_container
            test_nginx_config
            ;;
        "reload")
            log_info "重新加载 nginx 配置..."
            check_nginx_container
            reload_nginx
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
