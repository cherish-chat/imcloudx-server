// Code generated by goctl. DO NOT EDIT.
// Source: app.proto

package appservice

import (
	"context"

	"github.com/cherish-chat/imcloudx-server/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	VerifyClientReq  = pb.VerifyClientReq
	VerifyClientResp = pb.VerifyClientResp

	AppService interface {
		VerifyClient(ctx context.Context, in *VerifyClientReq, opts ...grpc.CallOption) (*VerifyClientResp, error)
	}

	defaultAppService struct {
		cli zrpc.Client
	}
)

func NewAppService(cli zrpc.Client) AppService {
	return &defaultAppService{
		cli: cli,
	}
}

func (m *defaultAppService) VerifyClient(ctx context.Context, in *VerifyClientReq, opts ...grpc.CallOption) (*VerifyClientResp, error) {
	client := pb.NewAppServiceClient(m.cli.Conn())
	return client.VerifyClient(ctx, in, opts...)
}
