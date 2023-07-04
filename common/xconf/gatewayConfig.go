package xconf

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type GatewayConfig struct {
	Port    int   `json:",default=12301"`
	Timeout int64 `json:",default=15000"`
	Http    struct {
		Cors struct {
			Enable           bool     `json:",optional"`
			AllowOrigins     []string `json:",optional"`
			AllowHeaders     []string `json:",optional"`
			AllowMethods     []string `json:",optional"`
			ExposeHeaders    []string `json:",optional"`
			AllowCredentials bool     `json:",optional"`
		} `json:",optional"`
		ApiLog struct {
			Apis []string `json:",optional"` // 格式: GET r'^/api/v1/user/.*' 表示所有以 /api/v1/user/ 开头的 GET 请求都会被记录
		}
		Host string `json:",default=0.0.0.0"`
		Port int    `json:",default=12300"`
	}
	Websocket struct {
		KeepAliveTickerSecond int `json:",default=30"` // 定时器，每隔n秒检测连接是否存活
		KeepAliveSecond       int `json:",default=60"` // 检测是否存活时，如果超过n秒没有收到客户端的消息，则关闭连接
	}
}

func (c Config) GatewayRpcServerConf() zrpc.RpcServerConf {
	name := "gateway"
	return zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: name,
			Log:  c.LogConf(name),
			Mode: c.Mode,
		},
		ListenOn: fmt.Sprintf(":%d", c.Gateway.Port),
		Etcd:     c.Etcd(name),
		Timeout:  c.Gateway.Timeout,
		Health:   true,
	}
}
