package svc

import (
	"github.com/cherish-chat/imcloudx-server/app/app/client/appservice"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	AppService appservice.AppService
}

func NewServiceContext(c config.Config) *ServiceContext {
	appClient := zrpc.MustNewClient(c.RpcClientConf.App)
	s := &ServiceContext{
		Config:     c,
		AppService: appservice.NewAppService(appClient),
	}
	return s
}
