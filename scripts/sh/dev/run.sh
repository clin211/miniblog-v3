#!/bin/bash

# 定义微服务列表
# 格式: "服务名:服务入口文件路径"
SERVICES=(
    "user-api:./apps/user/api/api.go:./apps/user/api/etc/api-api.yaml"
    "user-rpc:./apps/user/rpc/rpc.go:./apps/user/rpc/etc/rpc-api.yaml"
    # 如果有更多服务，继续在这里添加
    # "service-c:app/service-c/c.go:app/service-c/etc/service-c.yaml"
)


# 临时目录，用于存放编译后的二进制文件
TMP_DIR="./tmp"
# 创建临时目录
mkdir -p $TMP_DIR

# 在脚本退出时，杀掉所有由这个脚本启动的后台进程
# 添加了 -k 选项来确保命令在所有任务被杀死后立即返回
trap 'echo "Killing background processes..."; kill $(jobs -p)' EXIT

# 遍历服务列表，编译并启动
for service in "${SERVICES[@]}"; do
    # 解析服务名、入口路径和配置文件路径
    IFS=':' read -r service_name service_path config_path <<< "$service"

    echo "Building ${service_name}..."
    # 编译，将二进制文件输出到 tmp 目录下，以服务名命名
    go build -o "${TMP_DIR}/${service_name}" "${service_path}"

    if [ $? -eq 0 ]; then
        echo "Starting ${service_name} with config ${config_path}..."
        # 在后台启动服务, 并使用 -f 指定配置文件
        ./${TMP_DIR}/${service_name} -f "${config_path}" &
    else
        echo "Failed to build ${service_name}."
        # 如果任何一个服务编译失败，则退出脚本
        exit 1
    fi
done

echo "All services are running in the background. Watching for changes..."

# 等待所有后台任务完成
wait
