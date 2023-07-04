package handler

import (
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/handler/offerhandler"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {

	// ws
	{
		wsHandler := NewWsHandler(svcCtx)
		engine.GET("/ws", wsHandler.Upgrade)
	}

	// offer
	{
		offerHandler := offerhandler.NewOfferHandler(svcCtx)
		engine.GET("/offer", offerHandler.Offer)
		engine.POST("/offer", offerHandler.Offer)
	}
}
