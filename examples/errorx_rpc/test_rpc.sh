#!/bin/bash

# 测试 RPC 网关服务
echo "=== 测试 RPC 网关服务 ==="
echo "服务地址: http://localhost:8292"
echo ""

# 测试单个用户查询 - 成功案例
echo "1. 测试单个用户查询 - 成功案例"
curl -X GET "http://localhost:8292/rpc/users/123" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 测试单个用户查询 - 错误案例 (ID=0)
echo "2. 测试单个用户查询 - 错误案例 (ID=0)"
curl -X GET "http://localhost:8292/rpc/users/0" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 测试单个用户查询 - 错误案例 (ID=404)
echo "3. 测试单个用户查询 - 错误案例 (ID=404)"
curl -X GET "http://localhost:8292/rpc/users/404" \
  -H "Content-Type: application/json" \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 测试批量用户查询 - 混合成功和失败
echo "4. 测试批量用户查询 - 混合成功和失败"
curl -X POST "http://localhost:8292/rpc/users/batch" \
  -H "Content-Type: application/json" \
  -d '{"ids": ["123", "0", "404", "456"]}' \
  -w "\nHTTP状态码: %{http_code}\n\n"

# 测试批量用户查询 - 全部成功
echo "5. 测试批量用户查询 - 全部成功"
curl -X POST "http://localhost:8292/rpc/users/batch" \
  -H "Content-Type: application/json" \
  -d '{"ids": ["123", "456", "789"]}' \
  -w "\nHTTP状态码: %{http_code}\n\n"
