package main

import (
	"flag"
	"fmt"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/server"
	"github.com/cherish-chat/imcloudx-server/common/xconf"

	gatewayserviceServer "github.com/cherish-chat/imcloudx-server/app/gateway/internal/server/gatewayservice"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
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
	serverConf := c.GatewayRpcServerConf()
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(serverConf, func(grpcServer *grpc.Server) {
		pb.RegisterGatewayServiceServer(grpcServer, gatewayserviceServer.NewGatewayServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	httpServer := server.NewHttpServer(ctx)
	go httpServer.Start()

	fmt.Printf("Starting rpc server at %s...\n", serverConf.ListenOn)
	s.Start()
}
