package main

import (
	"flag"
	"fmt"

	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/config"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/handler"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/rpc-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
