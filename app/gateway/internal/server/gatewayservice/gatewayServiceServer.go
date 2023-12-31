// Code generated by goctl. DO NOT EDIT.
// Source: gateway.proto

package server

import (
	"context"

	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common/pb"
)

type GatewayServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedGatewayServiceServer
}

func NewGatewayServiceServer(svcCtx *svc.ServiceContext) *GatewayServiceServer {
	return &GatewayServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *GatewayServiceServer) GetSessionDescription(ctx context.Context, in *pb.GetSessionDescriptionReq) (*pb.GetSessionDescriptionResp, error) {
	l := gatewayservicelogic.NewGetSessionDescriptionLogic(ctx, s.svcCtx)
	return l.GetSessionDescription(in)
}

func (s *GatewayServiceServer) Node(ctx context.Context, in *pb.NodeReq) (*pb.NodeResp, error) {
	l := gatewayservicelogic.NewNodeLogic(ctx, s.svcCtx)
	return l.Node(in)
}
