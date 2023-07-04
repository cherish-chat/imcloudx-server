package handler

import (
	"context"
	"errors"
	gatewayservicelogic "github.com/cherish-chat/imcloudx-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common/pb"
	"github.com/cherish-chat/imcloudx-server/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"io"
	"math"
	"nhooyr.io/websocket"
	"strings"
)

type WsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewWsHandler(svcCtx *svc.ServiceContext) *WsHandler {
	return &WsHandler{
		svcCtx: svcCtx,
	}
}

func (h *WsHandler) Upgrade(ginCtx *gin.Context) {
	r := ginCtx.Request
	w := ginCtx.Writer
	logger := logx.WithContext(r.Context())
	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	compressionMode := websocket.CompressionNoContextTakeover
	// https://github.com/nhooyr/websocket/issues/218
	// 如果是Safari浏览器，不压缩
	if strings.Contains(r.UserAgent(), "Safari") {
		compressionMode = websocket.CompressionDisabled
	}
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:         nil,
		InsecureSkipVerify:   true,
		OriginPatterns:       nil,
		CompressionMode:      compressionMode,
		CompressionThreshold: 0,
	})
	if err != nil {
		// 如果是 / 说明是健康检查
		if r.URL.Path == "/" {
			return
		}
		logger.Errorf("failed to accept websocket connection: %v", err)
		return
	}
	c.SetReadLimit(math.MaxInt32)
	verifyClientResp, err := h.svcCtx.AppService.VerifyClient(r.Context(), &pb.VerifyClientReq{
		ClientId:     r.URL.Query().Get("clientId"),
		ClientSecret: r.URL.Query().Get("clientSecret"),
	})
	if err != nil {
		logger.Errorf("beforeConnect error: %v", err)
		c.Close(websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeServerInternalError), err.Error())
		return
	}
	if !verifyClientResp.Ok {
		c.Close(websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeAuthenticationFailed), verifyClientResp.Tip)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	ctx, cancelFunc := context.WithCancel(r.Context())
	connectionId := utils.Snowflake.Int64()
	defer func() {
		logger.Debugf("removing subscriber: %d", connectionId)
		err := gatewayservicelogic.WsManager.RemoveSubscriber(verifyClientResp.AppId, connectionId, websocket.StatusNormalClosure, "finished")
		if err != nil {
			logger.Errorf("failed to remove subscriber: %v", err)
			return
		} else {
			logger.Debugf("removed subscriber: %d", connectionId)
		}
	}()
	connection, err := gatewayservicelogic.WsManager.AddSubscriber(ctx, verifyClientResp.AppId, c, connectionId)
	if err != nil {
		logger.Errorf("failed to add subscriber: %v", err)
		c.Close(websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeServerInternalError), err.Error())
		cancelFunc()
		return
	}
	go func() {
		// 读取消息
		defer cancelFunc()
		for {
			logger.Debugf("start read")
			typ, msg, err := c.Read(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					// 正常关闭
				} else if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
					websocket.CloseStatus(err) == websocket.StatusGoingAway {
					// 正常关闭
					logx.Infof("websocket closed: %v", err)
				} else if strings.Contains(err.Error(), "connection reset by peer") {
					// 网络断开
					logx.Infof("websocket closed: %v", err)
				} else if strings.Contains(err.Error(), "corrupt input") {
					// 输入数据错误
					logx.Infof("websocket closed: %v", err)
				} else {
					logx.Errorf("failed to read message: %v", err)
				}
				return
			}
			go func() {
				_ = h.onReceive(ctx, connection, typ, msg)
			}()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

type WsReceiveMessageBinary = pb.WsReceiveMessageBinary

func (h *WsHandler) onReceive(ctx context.Context, connection *gatewayservicelogic.WsConnection, typ websocket.MessageType, msg []byte) error {
	if typ == websocket.MessageText {
		gatewayservicelogic.WsManager.KeepAlive(ctx, connection)
		return nil
	} else if typ == websocket.MessageBinary {
		message := &WsReceiveMessageBinary{}
		err := proto.Unmarshal(msg, message)
		if err != nil {
			logx.Errorf("failed to unmarshal message: %v", err)
			return err
		}
		switch message.Type {
		case pb.WsReceiveMessageBinaryType_ReceiveRequest:
		// TODO onReceiveRequest
		case pb.WsReceiveMessageBinaryType_ReceiveResponse:
			response := &pb.NodeResp{}
			err := proto.Unmarshal(message.Data, response)
			if err != nil {
				logx.Errorf("failed to unmarshal response: %v", err)
				return err
			}
			gatewayservicelogic.WsManager.OnReceiveResponse(ctx, connection, response)
		}
	}
	return nil
}
