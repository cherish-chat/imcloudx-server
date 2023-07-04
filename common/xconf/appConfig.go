package xconf

import (
	"fmt"
	"github.com/cherish-chat/imcloudx-server/app/app/client/appservice"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"time"
)

type AppConfig struct {
	Port    int   `json:",default=12301"`
	Timeout int64 `json:",default=15000"`
}

func (c Config) AppRpcServerConf() zrpc.RpcServerConf {
	name := "app"
	return zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: name,
			Log:  c.LogConf(name),
			Mode: c.Mode,
		},
		ListenOn: fmt.Sprintf(":%d", c.App.Port),
		Etcd:     c.Etcd(name),
		Timeout:  c.App.Timeout,
		Health:   true,
	}
}

func (c Config) NewAppRpc() appservice.AppService {
	var clientConf zrpc.RpcClientConf
	rpcConfig := c.App
	name := "app"
	// 判断是否启用的etcd服务发现
	if len(c.Discovery.Etcd.Hosts) > 0 {
		// 启用etcd服务发现
		clientConf = zrpc.RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts:              c.Discovery.Etcd.Hosts,
				Key:                c.Discovery.Etcd.KeyPrefix + name,
				ID:                 c.Discovery.Etcd.ID,
				User:               c.Discovery.Etcd.User,
				Pass:               c.Discovery.Etcd.Pass,
				CertFile:           c.Discovery.Etcd.CertFile,
				CertKeyFile:        c.Discovery.Etcd.CertKeyFile,
				CACertFile:         c.Discovery.Etcd.CACertFile,
				InsecureSkipVerify: c.Discovery.Etcd.InsecureSkipVerify,
			},
			Endpoints: nil,
			Target:    "",
			NonBlock:  true,
			Timeout:   rpcConfig.Timeout,
		}
	} else if c.Discovery.K8sNamespace != "" {
		// 启用k8s服务发现
		clientConf = zrpc.RpcClientConf{
			Endpoints: nil,
			Target:    fmt.Sprintf("k8s://%s/%s-svc:%d", c.Discovery.K8sNamespace, name, rpcConfig.Port),
			NonBlock:  true,
			Timeout:   rpcConfig.Timeout,
		}
	} else {
		// 不启用服务发现
		clientConf = zrpc.RpcClientConf{
			Endpoints: []string{
				fmt.Sprintf("localhost:%d", rpcConfig.Port),
			},
			NonBlock: true,
			Timeout:  rpcConfig.Timeout,
		}
	}
	client, err := zrpc.NewClient(clientConf, zrpc.WithTimeout(time.Millisecond*time.Duration(rpcConfig.Timeout)))
	if err != nil {
		logx.Errorf("new user rpc client failed: %s", err.Error())
		os.Exit(1)
		return nil
	}
	return appservice.NewAppService(client)
}
