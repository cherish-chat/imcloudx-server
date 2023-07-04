package svc

import (
	"github.com/cherish-chat/imcloudx-server/common/xconf"
)

type ServiceContext struct {
	Config xconf.Config
}

func NewServiceContext(c xconf.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
