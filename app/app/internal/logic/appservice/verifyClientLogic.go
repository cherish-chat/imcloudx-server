package appservicelogic

import (
	"context"

	"github.com/cherish-chat/imcloudx-server/app/app/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyClientLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyClientLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyClientLogic {
	return &VerifyClientLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyClientLogic) VerifyClient(in *pb.VerifyClientReq) (*pb.VerifyClientResp, error) {
	// todo: add your logic here and delete this line

	return &pb.VerifyClientResp{Ok: true, AppId: "1"}, nil
}
