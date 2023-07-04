package gatewayservicelogic

import (
	"context"
	"math/rand"
	"time"

	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type NodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NodeLogic {
	return &NodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NodeLogic) Node(in *pb.NodeReq) (*pb.NodeResp, error) {
	// TODO 获取所有pod中的connection
	logx.Infof("node request: %+v", in)
	connections, ok := WsManager.wsConnectionMap.GetByAppId(in.AppId)
	if !ok || len(connections) == 0 {
		return &pb.NodeResp{
			AppId:   in.AppId,
			Headers: nil,
			Status:  502,
			Body:    nil,
			ErrMsg:  "node not found",
		}, nil
	}
	// 随机选择一个节点
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(connections))
	connection := connections[index]
	response, err := WsManager.SendRequest(connection, in, time.Second*5)
	if err != nil {
		return &pb.NodeResp{
			AppId:   in.AppId,
			Headers: nil,
			Status:  500,
			Body:    nil,
			ErrMsg:  err.Error(),
		}, nil
	}
	return response, nil
}
