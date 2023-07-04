package gatewayservicelogic

import (
	"context"

	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionDescriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionDescriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionDescriptionLogic {
	return &GetSessionDescriptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionDescriptionLogic) GetSessionDescription(in *pb.GetSessionDescriptionReq) (*pb.GetSessionDescriptionResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetSessionDescriptionResp{}, nil
}
