#!/bin/bash

echo "=== grpcurl 完整测试 ==="
echo "服务地址: localhost:50051"
echo ""

# 检查 grpcurl 是否安装
if ! command -v grpcurl &> /dev/null; then
    echo "grpcurl 未安装，正在安装..."
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
fi

echo "1. 列出所有服务"
grpcurl -plaintext localhost:50051 list
echo ""

echo "2. 列出 UserService 的所有方法"
grpcurl -plaintext localhost:50051 list user.UserService
echo ""

echo "3. 查看 GetUser 方法的详细信息"
grpcurl -plaintext localhost:50051 describe user.UserService.GetUser
echo ""

echo "4. 测试 GetUser - 成功案例"
grpcurl -plaintext -d '{"id": "123"}' localhost:50051 user.UserService.GetUser
echo ""

echo "5. 测试 GetUser - 错误案例 (ID=0)"
grpcurl -plaintext -d '{"id": "0"}' localhost:50051 user.UserService.GetUser
echo ""

echo "6. 测试 GetUser - 错误案例 (ID=404)"
grpcurl -plaintext -d '{"id": "404"}' localhost:50051 user.UserService.GetUser
echo ""

echo "7. 测试 BatchGetUsers - 混合成功和失败"
grpcurl -plaintext -d '{"ids": ["123", "0", "404", "456"]}' localhost:50051 user.UserService.BatchGetUsers
echo ""

echo "8. 测试 BatchGetUsers - 全部成功"
grpcurl -plaintext -d '{"ids": ["123", "456", "789"]}' localhost:50051 user.UserService.BatchGetUsers
echo ""

echo "9. 使用详细输出模式"
grpcurl -v -plaintext -d '{"id": "123"}' localhost:50051 user.UserService.GetUser
echo ""

echo "=== 测试完成 ==="
