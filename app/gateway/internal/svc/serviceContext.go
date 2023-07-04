package svc

import (
	"github.com/cherish-chat/imcloudx-server/app/app/client/appservice"
	"github.com/cherish-chat/imcloudx-server/common/xconf"
)

type ServiceContext struct {
	Config     xconf.Config
	AppService appservice.AppService
}

func NewServiceContext(c xconf.Config) *ServiceContext {
	s := &ServiceContext{
		Config:     c,
		AppService: c.NewAppRpc(),
	}
	return s
}
