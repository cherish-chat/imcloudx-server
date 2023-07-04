package offerhandler

import (
	gatewayservicelogic "github.com/cherish-chat/imcloudx-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common/pb"
	"github.com/cherish-chat/imcloudx-server/common/utils"
	"github.com/gin-gonic/gin"
)

type OfferHandler struct {
	svcCtx *svc.ServiceContext
}

func NewOfferHandler(svcCtx *svc.ServiceContext) *OfferHandler {
	return &OfferHandler{svcCtx: svcCtx}
}

type OfferReq struct {
	AppId string     `json:"appId" form:"appId"`
	Sdp   string     `json:"sdp" form:"sdp"`
	Type  pb.SDPType `json:"type" form:"type"`
}

type Answer struct {
	Sdp  string     `json:"sdp" form:"sdp"`
	Type pb.SDPType `json:"type" form:"type"`
}

func (h *OfferHandler) Offer(ctx *gin.Context) {
	var req OfferReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}
	nodeResp, err := gatewayservicelogic.NewNodeLogic(ctx.Request.Context(), h.svcCtx).Node(&pb.NodeReq{
		AppId:   req.AppId,
		Headers: map[string]string{},
		Method:  "/offer",
		Body:    utils.Json.MarshalToBytes(req),
	})
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}
	if nodeResp.Status != 200 {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  nodeResp.ErrMsg,
		})
		return
	}
	answer := &Answer{}
	if err := utils.Json.Unmarshal(nodeResp.Body, answer); err != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  "",
			"data": answer,
		})
		return
	}
}
