package svc

import (
	"github.com/clin211/miniblog-v3/apps/user/models"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel models.UsersModel
	// 添加原始数据库连接用于事务处理
	DB        sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 连接 MySQL 数据库
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// 初始化用户模型
	userModel := models.NewUsersModel(conn, c.Cache)

	return &ServiceContext{
		Config:    c,
		UserModel: userModel,
		DB:        conn, // 保存原始连接
	}
}
