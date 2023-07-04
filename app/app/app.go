package main

import (
	"flag"
	"fmt"
	"github.com/cherish-chat/imcloudx-server/common/xconf"

	appserviceServer "github.com/cherish-chat/imcloudx-server/app/app/internal/server/appservice"
	"github.com/cherish-chat/imcloudx-server/app/app/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c xconf.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	serverConf := c.AppRpcServerConf()
	s := zrpc.MustNewServer(serverConf, func(grpcServer *grpc.Server) {
		pb.RegisterAppServiceServer(grpcServer, appserviceServer.NewAppServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", serverConf.ListenOn)
	s.Start()
}
