package main

import (
	"flag"
	"fmt"

	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/config"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/server"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"
	"github.com/clin211/miniblog-v3/pkg/middleware"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	logx.MustSetup(logx.LogConf{
		ServiceName:      "user-rpc", // 服务名称
		Mode:             "file",     // 日志模式
		Path:             "./logs",   // 日志文件存储路径
		Level:            "info",     // 日志级别
		MaxSize:          100,        // 每个日志文件的最大大小，单位MB
		MaxContentLength: 200,        // 日志长度限制
		MaxBackups:       10,         // 文件输出模式，按照大小分割时，最多文件保留个数
		Compress:         true,       // 是否启用日志压缩
	})

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rpc.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		fmt.Println("RegisterUserServer", c.Mode)
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 添加gRPC拦截器
	s.AddUnaryInterceptors(middleware.AuthnInterceptor())

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
