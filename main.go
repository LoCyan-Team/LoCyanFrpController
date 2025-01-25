package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"lcf-controller/inject"
	"lcf-controller/logger"
	"lcf-controller/net/server"
	"lcf-controller/pkg/config"
	"lcf-controller/pkg/info"
	_type "lcf-controller/pkg/type/frps"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewWebSocket 初始化WebSocket客户端
func NewWebSocket() *WsClient {
	ws := new(WsClient)
	cfgInfo := config.ReadCfg()
	ws.addr = cfgInfo.Addr
	ws.config = cfgInfo
	return ws
}

// ConnectWsServer 连接到WebSocket服务器
func (w *WsClient) ConnectWsServer() (err error) {
	conn, _, err := websocket.DefaultDialer.Dial(w.addr, nil)
	if err != nil {
		return err
	}
	w.conn = conn
	return nil
}

// SendMsg 发送消息到服务器
func (w *WsClient) SendMsg(action string, data map[string]any) (err error) {
	req := new(BasicRequest)
	req.Action = action
	req.Node.Id = w.config.NodeId
	req.Node.ApiKey = w.config.NodeApiKey
	req.Data = data
	msg, err := json.Marshal(req)
	err = w.conn.WriteMessage(websocket.TextMessage, msg)

	if err != nil {
		logger.Logger.Fatal("failed to send message", zap.Error(err))
	}

	return nil
}

// ReadMsg 从服务器读取消息
func (w *WsClient) ReadMsg() {
	defer func() {
		err := w.conn.Close()
		if err != nil {
			logger.Logger.Error("error closing connection", zap.Error(err))
		}
	}()

	for {
		_, msg, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Logger.Error("error reading message", zap.Error(err))
			}
			break
		}
		var msgJson WsResponse
		err = json.Unmarshal(msg, &msgJson)
		if err != nil {
			logger.Logger.Error("can't unmarshal json message", zap.Error(err))
		}
		if msgJson.Status != 200 {
			logger.Logger.Error("error Message from server", zap.String("msg", string(msg)))
		}
		if msgJson.Status == 200 {
			logger.Logger.Debug("Received message from server", zap.String("msg", string(msg)))
		}
	}
}

func (w *WsClient) sendNodeStatsToServer(serverInfo _type.ServerInfoResponse) {
	// nodeInfo
	err := w.SendMsg("upload-node-stats", info.GetNodeInfo(serverInfo))
	if err != nil {
		logger.Logger.Error("send node info to server failed!", zap.Error(err))
	}
}

func (w *WsClient) sendProxyStatsToServer() {
	types := []string{"tcp", "udp", "http", "https", "xtcp", "stcp"}
	for _, p := range types {
		proxies := info.GetProxies(p)
		for _, j := range proxies {
			err := w.SendMsg("upload-proxy-stats", j)
			logger.Logger.Info("send proxy info to the server")
			if err != nil {
				logger.Logger.Error("send proxy info to server failed!", zap.Error(err))
			}
		}
	}
}

// WsClient WebSocket客户端结构
type WsClient struct {
	addr   string
	conn   *websocket.Conn
	config *config.Config
}

type WsResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type BasicRequest struct {
	Action string         `json:"action"`
	Node   NodeInfo       `json:"node"`
	Data   map[string]any `json:"data"`
}

type NodeInfo struct {
	Id     int    `json:"id"`
	ApiKey string `json:"api_key"`
}

func createContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Graceful shutdown
		shutdownChan := make(chan os.Signal, 1)
		signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
		<-shutdownChan
		logger.Logger.Info("shutting down gracefully...")

		logger.Logger.Info("closing OpenGFW engine...")
		cancel()
		logger.Logger.Info("OpenGFW engine closed")

		os.Exit(0)
	}()
	return ctx, cancel
}

func main() {
	logger.InitLogger()
	ctx, _ := createContext()

	go inject.RunOpenGFW(ctx)

	ws := NewWebSocket()
	logger.Logger.Info("starting to connect WebSocket...")
	err := ws.ConnectWsServer()
	if err != nil {
		logger.Logger.Fatal(
			"can't connect to WebSocket server",
			zap.Error(err),
		)
	} else {
		logger.Logger.Info("connect to WebSocket server successfully!")
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				logger.Logger.Fatal(
					"can't close WebSocket connection",
					zap.Error(err),
				)
			}
		}(ws.conn)
		go ws.ReadMsg()
		ticker := time.NewTicker(time.Second * ws.config.SendDuration)
		defer ticker.Stop()

		serverInfo := server.GetServerInfo()
		ws.sendNodeStatsToServer(serverInfo)
		ws.sendProxyStatsToServer()

		for range ticker.C {
			serverInfo := server.GetServerInfo()
			ws.sendNodeStatsToServer(serverInfo)
			ws.sendProxyStatsToServer()
		}
	}
}
