package gatewayservicelogic

import (
	"context"
	"errors"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common/pb"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type WsConnectionMap struct {
	idConnectionMap     map[int64]*WsConnection
	idConnectionMapLock sync.RWMutex
	appIdsMap           map[string][]*WsConnection
	appIdsMapLock       sync.RWMutex
	idAliveTimeMap      map[int64]time.Time
	idAliveTimeMapLock  sync.RWMutex

	responseChanMap sync.Map // key = requestId, value = chan *pb.ResponseHeader
}

var ErrTimeout = errors.New("timeout")

func (w *WsConnectionMap) waitResponse(requestId string, timeout time.Duration) (*pb.NodeResp, error) {
	ch := make(chan *pb.NodeResp, 0)
	w.responseChanMap.Store(requestId, ch)
	defer w.responseChanMap.Delete(requestId)
	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

func (w *WsConnectionMap) onResponse(response *pb.NodeResp) {
	requestId := response.RequestId
	if ch, ok := w.responseChanMap.Load(requestId); ok {
		ch.(chan *pb.NodeResp) <- response
	}
}

func (w *WsConnectionMap) GetAll() []*WsConnection {
	var connections []*WsConnection
	w.idConnectionMapLock.RLock()
	defer w.idConnectionMapLock.RUnlock()
	for _, v := range w.idConnectionMap {
		connections = append(connections, v)
	}
	return connections
}

func (w *WsConnectionMap) GetByAppId(appId string) ([]*WsConnection, bool) {
	w.appIdsMapLock.RLock()
	defer w.appIdsMapLock.RUnlock()
	v, ok := w.appIdsMap[appId]
	return v, ok
}

func (w *WsConnectionMap) Delete(userId string, connectionId int64) {
	w.idConnectionMapLock.Lock()
	delete(w.idConnectionMap, connectionId)
	w.idConnectionMapLock.Unlock()
	w.appIdsMapLock.Lock()
	defer w.appIdsMapLock.Unlock()
	//获取用户的所有连接
	connections, ok := w.appIdsMap[userId]
	if !ok {
		return
	}
	//删除用户的某个连接
	var newConnections []*WsConnection
	for _, connection := range connections {
		if connection.Id != connectionId {
			newConnections = append(newConnections, connection)
		}
	}
	w.appIdsMap[userId] = newConnections
	w.idAliveTimeMapLock.Lock()
	delete(w.idAliveTimeMap, connectionId)
	w.idAliveTimeMapLock.Unlock()
}

func (w *WsConnectionMap) GetByConnectionId(connectionId int64) (*WsConnection, bool) {
	// RLock() 读锁
	w.idConnectionMapLock.RLock()
	defer w.idConnectionMapLock.RUnlock()
	v, ok := w.idConnectionMap[connectionId]
	return v, ok
}

func (w *WsConnectionMap) GetAliveTime(connectionId int64) (time.Time, bool) {
	// RLock() 读锁
	w.idAliveTimeMapLock.RLock()
	defer w.idAliveTimeMapLock.RUnlock()
	v, ok := w.idAliveTimeMap[connectionId]
	return v, ok
}

type wsManager struct {
	svcCtx          *svc.ServiceContext
	wsConnectionMap *WsConnectionMap
}

var WsManager *wsManager

type WsConnection struct {
	Id         int64
	Connection *websocket.Conn
	AppId      string
	Ctx        context.Context
}

func InitWsManager(svcCtx *svc.ServiceContext) {
	WsManager = &wsManager{
		svcCtx: svcCtx,
		wsConnectionMap: &WsConnectionMap{
			appIdsMap:           make(map[string][]*WsConnection),
			appIdsMapLock:       sync.RWMutex{},
			idConnectionMap:     make(map[int64]*WsConnection),
			idConnectionMapLock: sync.RWMutex{},
			idAliveTimeMap:      make(map[int64]time.Time),
			idAliveTimeMapLock:  sync.RWMutex{},
		},
	}
	go WsManager.loopCheck()
}

func (w *wsManager) loopCheck() {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(w.svcCtx.Config.Websocket.KeepAliveTickerSecond))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				connections := w.wsConnectionMap.GetAll()
				for _, connection := range connections {
					_, ok := w.wsConnectionMap.GetAliveTime(connection.Id)
					if !ok {
						// 删除连接
						w.RemoveSubscriber(connection.AppId, connection.Id, websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeHeartbeatTimeout), "heartbeat timeout")
					}
				}
			}
		}
	}()
}
func (w *wsManager) RemoveSubscriber(appId string, id int64, closeCode websocket.StatusCode, closeReason string) error {
	connection, ok := w.wsConnectionMap.GetByConnectionId(id)
	if ok {
		_ = connection.Connection.Close(closeCode, closeReason)
	}
	w.wsConnectionMap.Delete(appId, id)
	go func() {
		if ok {
			//_, e := w.svcCtx.CallbackService.AppAfterOffline(context.Background(), &pb.UserAfterOfflineReq{Header: header})
			//if e != nil {
			//	logx.Errorf("UserAfterOffline error: %s", e.Error())
			//}
		}
	}()
	return nil
}

// clearConnectionTimer 定时器清除连接
func (w *wsManager) clearConnectionTimer(connection *WsConnection) {
	ticker := time.NewTicker(time.Second * time.Duration(w.svcCtx.Config.Websocket.KeepAliveTickerSecond))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			//使用 id 查询连接最后活跃时间
			aliveTime, ok := w.wsConnectionMap.GetAliveTime(connection.Id)
			sub := time.Now().Sub(aliveTime)
			if !ok || sub > time.Second*time.Duration(w.svcCtx.Config.Websocket.KeepAliveSecond) {
				// 删除连接
				w.RemoveSubscriber(connection.AppId, connection.Id, websocket.StatusCode(pb.WebsocketCustomCloseCode_CloseCodeHeartbeatTimeout), "heartbeat timeout")
				return
			}
		}
	}
}

func (w *WsConnectionMap) SetAliveTime(ctx context.Context, connectionId int64, aliveTime time.Time) {
	w.idAliveTimeMapLock.Lock()
	w.idAliveTimeMap[connectionId] = aliveTime
	w.idAliveTimeMapLock.Unlock()
}

func (w *WsConnectionMap) Set(connectionId int64, value *WsConnection) {
	w.idConnectionMapLock.Lock()
	w.idConnectionMap[connectionId] = value
	w.idConnectionMapLock.Unlock()
	w.appIdsMapLock.Lock()
	w.appIdsMap[value.AppId] = append(w.appIdsMap[value.AppId], value)
	w.appIdsMapLock.Unlock()
	w.idAliveTimeMapLock.Lock()
	w.idAliveTimeMap[connectionId] = time.Now()
	w.idAliveTimeMapLock.Unlock()
}

func (w *wsManager) KeepAlive(ctx context.Context, connection *WsConnection) {
	w.wsConnectionMap.SetAliveTime(ctx, connection.Id, time.Now())
}

func (w *wsManager) AddSubscriber(ctx context.Context, appId string, connection *websocket.Conn, id int64) (*WsConnection, error) {
	wsConnection := &WsConnection{
		Id:         id,
		Connection: connection,
		AppId:      appId,
		Ctx:        ctx,
	}
	//启动定时器 定时删掉连接
	go w.clearConnectionTimer(wsConnection)
	w.wsConnectionMap.Set(id, wsConnection)
	go func() {
		//_, e := w.svcCtx.CallbackService.UserAfterOnline(ctx, &pb.UserAfterOnlineReq{Header: header})
		//if e != nil {
		//	logx.Errorf("UserAfterOnline error: %s", e.Error())
		//}
	}()
	return wsConnection, nil
}

func (w *wsManager) SendRequest(connection *WsConnection, in *pb.NodeReq, timeout time.Duration) (*pb.NodeResp, error) {
	response, err := w.wsConnectionMap.waitResponse(in.RequestId, timeout)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (w *wsManager) OnReceiveResponse(ctx context.Context, connection *WsConnection, response *pb.NodeResp) {
	w.wsConnectionMap.onResponse(response)
}
