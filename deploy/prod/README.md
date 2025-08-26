# 生产环境部署指南

本文档描述了 miniblog-v3 项目的生产环境部署流程，采用基础设施与应用服务分离的架构设计。

## 架构设计

### 服务分离架构

- **基础设施服务**: MySQL、Redis、etcd、Zookeeper、Kafka
- **应用服务**: user-api、user-rpc 等业务服务
- **外部服务**: 使用服务器上现有的 nginx 容器

### 优势

- ✅ 基础设施服务独立管理，更新频率低
- ✅ 应用服务可以频繁更新而不影响基础设施
- ✅ 数据持久化，避免意外丢失
- ✅ 更好的资源隔离和管理

## 目录结构

```
deploy/prod/
├── docker-compose.yml                    # 应用服务配置
├── docker-compose.env.yml                # 基础设施服务配置（包含环境变量）
├── infrastructure-manager.sh             # 基础设施管理脚本
├── kafka-manager.sh                      # Kafka 管理脚本
├── deploy.sh                             # 应用服务部署脚本
├── env.example                           # 环境变量示例
├── nginx/                                # Nginx 配置目录
│   ├── conf.d/
│   │   ├── miniblog.conf                 # 主配置文件
│   │   └── services/                     # 服务配置目录
│   │       ├── user.conf                 # 用户服务配置
│   │       └── template.conf             # 配置模板
│   └── nginx-manager.sh                  # Nginx 管理脚本
└── README.md                             # 本文档
```

## 部署流程

### 第一步：启动基础设施服务

```bash
# 1. 启动基础设施服务
./infrastructure-manager.sh start

# 2. 检查服务状态
./infrastructure-manager.sh status

# 3. 创建 Kafka 默认主题（可选）
./kafka-manager.sh create-default
```

### 第二步：部署应用服务

```bash
# 1. 部署应用服务
./deploy.sh

# 2. 检查应用服务状态
docker-compose ps
```

### 第三步：配置外部 Nginx

```bash
# 1. 安装 nginx 配置
./nginx-manager.sh install

# 2. 检查 nginx 配置
./nginx-manager.sh test
```

## 服务管理

### 基础设施服务管理

```bash
# 启动基础设施服务
./infrastructure-manager.sh start

# 停止基础设施服务
./infrastructure-manager.sh stop

# 重启基础设施服务
./infrastructure-manager.sh restart

# 查看服务状态
./infrastructure-manager.sh status

# 查看服务日志
./infrastructure-manager.sh logs [service]

# 备份数据库
./infrastructure-manager.sh backup [backup-dir]

# 恢复数据库
./infrastructure-manager.sh restore <backup-file>

# 清理所有数据（危险操作）
./infrastructure-manager.sh clean
```

### Kafka 管理

```bash
# 列出所有主题
./kafka-manager.sh list

# 创建主题
./kafka-manager.sh create <topic-name> [partitions] [replicas]

# 查看主题详情
./kafka-manager.sh describe <topic-name>

# 发送测试消息
./kafka-manager.sh send <topic-name> [message]

# 消费消息
./kafka-manager.sh consume <topic-name> [group] [offset]

# 查看消费者组
./kafka-manager.sh list-groups

# 创建默认主题
./kafka-manager.sh create-default
```

### 应用服务管理

```bash
# 部署应用服务
./deploy.sh

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f [service]

# 重启服务
docker-compose restart [service]

# 停止服务
docker-compose down
```

### Nginx 服务管理

```bash
# 安装 nginx 配置
./nginx-manager.sh install

# 测试 nginx 配置
./nginx-manager.sh test

# 重新加载 nginx 配置
./nginx-manager.sh reload

# 恢复 nginx 配置
./nginx-manager.sh restore

# 手动添加新服务
# 在 nginx/conf.d/services/ 目录下创建对应的 .conf 文件
```

## 环境变量配置

### 数据库配置

```bash
MYSQL_ROOT_PASSWORD=root123456
MYSQL_DATABASE=miniblog
MYSQL_USER=miniblog
MYSQL_PASSWORD=miniblog123
```

### Redis 配置

```bash
REDIS_PASSWORD=redis123
```

### etcd 配置

```bash
ETCD_ENDPOINTS=http://localhost:2379
ETCD_CLIENT_PORT=2379
ETCD_PEER_PORT=2380
ETCD_MAX_REQUEST_BYTES=1048576
ETCD_QUOTA_BACKEND_BYTES=2147483648
```

### Zookeeper 配置

```bash
ZOOKEEPER_PORT=2181
ZOOKEEPER_TICK_TIME=2000
ZOOKEEPER_INIT_LIMIT=5
ZOOKEEPER_SYNC_LIMIT=2
ZOOKEEPER_MAX_CLIENT_CNXNS=60
ZOOKEEPER_AUTOPURGE_SNAP_RETAIN_COUNT=3
ZOOKEEPER_AUTOPURGE_PURGE_INTERVAL=24
```

### Kafka 配置 (Bitnami 镜像)

```bash
KAFKA_BOOTSTRAP_SERVERS=localhost:29092
KAFKA_ZOOKEEPER_CONNECT=localhost:2181
KAFKA_PORT=29092
KAFKA_BROKER_ID=1
KAFKA_JMX_PORT=9101
KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1
KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1
KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1
KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0
KAFKA_AUTO_CREATE_TOPICS_ENABLE=true
KAFKA_DELETE_TOPIC_ENABLE=true
KAFKA_LOG_RETENTION_HOURS=168
KAFKA_LOG_RETENTION_BYTES=1073741824
KAFKA_LOG_SEGMENT_BYTES=1073741824
KAFKA_NUM_PARTITIONS=3
KAFKA_DEFAULT_REPLICATION_FACTOR=1
KAFKA_MIN_INSYNC_REPLICAS=1
```

**注意**: 使用 Bitnami Kafka 镜像，所有配置参数都以 `KAFKA_CFG_` 前缀开头。

### Docker 镜像配置

```bash
DOCKER_REGISTRY=clin211
IMAGE_TAG=latest
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| MySQL | 3306 | 数据库服务 |
| Redis | 6379 | 缓存服务 |
| etcd | 2379 | 配置中心 |
| Zookeeper | 2181 | Kafka 协调服务 |
| Kafka | 29092 | 消息队列 |
| User API | 8888 | 用户 API 服务 |
| User RPC | 8080 | 用户 RPC 服务 |

## 健康检查

### 基础设施服务健康检查

```bash
# 检查所有基础设施服务
./infrastructure-manager.sh status

# 检查特定服务
docker exec miniblog-mysql mysqladmin ping -h localhost
docker exec miniblog-redis redis-cli ping
docker exec miniblog-etcd etcdctl endpoint health
docker exec miniblog-kafka kafka-topics --bootstrap-server localhost:9092 --list
```

### 应用服务健康检查

```bash
# 检查应用服务
curl http://localhost:8888/health
curl http://localhost:8080/health

# 检查 nginx 配置
curl https://your-domain.com/health
```

## 日志管理

### 日志文件位置

- **应用服务日志**: `./logs/user-api/`, `./logs/user-rpc/`
- **基础设施日志**: 通过 `./infrastructure-manager.sh logs` 查看
- **Nginx 日志**: 在外部 nginx 容器中查看

### 日志查看命令

```bash
# 查看应用服务日志
docker-compose logs -f [service]

# 查看基础设施服务日志
./infrastructure-manager.sh logs [service]

# 查看 nginx 日志
docker exec nginx tail -f /var/log/nginx/access.log
```

## 备份和恢复

### 数据库备份

```bash
# 创建备份
./infrastructure-manager.sh backup ./backups

# 恢复备份
./infrastructure-manager.sh restore ./backups/mysql_backup_20241201_120000.sql
```

### 数据卷备份

```bash
# 备份数据卷
docker run --rm -v miniblog-infrastructure_mysql_data:/data -v $(pwd):/backup alpine tar czf /backup/mysql_data_backup.tar.gz -C /data .

# 恢复数据卷
docker run --rm -v miniblog-infrastructure_mysql_data:/data -v $(pwd):/backup alpine tar xzf /backup/mysql_data_backup.tar.gz -C /data
```

## 故障排除

### 常见问题

1. **基础设施服务启动失败**

   ```bash
   # 检查 Docker 是否运行
   docker info
   
   # 检查端口占用
   netstat -tlnp | grep :3306
   
   # 查看详细日志
   ./infrastructure-manager.sh logs
   ```

2. **应用服务无法连接基础设施**

   ```bash
   # 检查网络连接
   docker network ls
   docker network inspect miniblog-network
   
   # 检查服务状态
   ./infrastructure-manager.sh status
   ```

3. **Kafka 主题创建失败**

   ```bash
   # 检查 Kafka 状态
   ./kafka-manager.sh cluster
   
   # 重新创建主题
   ./kafka-manager.sh create <topic-name>
   ```

4. **Nginx 配置问题**

   ```bash
   # 测试配置
   ./nginx-manager.sh test
   
   # 查看 nginx 错误日志
   docker exec nginx tail -f /var/log/nginx/error.log
   ```

### 性能优化

1. **数据库优化**

   ```bash
   # 调整 MySQL 配置
   docker exec miniblog-mysql mysql -e "SHOW VARIABLES LIKE 'max_connections';"
   ```

2. **Redis 优化**

   ```bash
   # 查看 Redis 内存使用
   docker exec miniblog-redis redis-cli info memory
   ```

3. **Kafka 优化**

   ```bash
   # 查看 Kafka 配置
   ./kafka-manager.sh cluster
   ```

## 监控和维护

### 资源监控

```bash
# 查看容器资源使用
docker stats

# 查看磁盘使用
df -h

# 查看内存使用
free -h
```

### 定期维护

1. **日志轮转**: 配置 logrotate 定期清理日志
2. **数据备份**: 定期备份数据库和重要数据
3. **安全更新**: 定期更新基础镜像和安全补丁
4. **性能监控**: 监控服务响应时间和资源使用

## CI/CD 集成

### GitHub Actions 部署

1. 推送代码到 `release` 分支触发自动部署
2. 自动构建和推送 Docker 镜像
3. 自动部署到服务器
4. 自动配置外部 nginx

### 手动部署

```bash
# 在服务器上执行
cd /path/to/miniblog-v3

# 启动基础设施（如果未启动）
./infrastructure-manager.sh start

# 部署应用服务
./deploy.sh

# 配置 nginx
./nginx-manager.sh install
```

## 安全建议

1. **修改默认密码**: 修改所有服务的默认密码
2. **网络安全**: 配置防火墙，只开放必要端口
3. **SSL 证书**: 使用有效的 SSL 证书
4. **访问控制**: 限制数据库和服务的访问权限
5. **定期备份**: 建立自动备份机制

## 联系支持

如果遇到问题，请：

1. 查看相关日志文件
2. 检查服务状态
3. 参考故障排除部分
4. 提交 Issue 到项目仓库

---

通过这种分离的架构设计，你可以更灵活地管理基础设施和应用服务，提高系统的可维护性和稳定性。
