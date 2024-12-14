package main

import (
	"LoCyanFrpController/net/server"
	"LoCyanFrpController/pkg/config"
	"LoCyanFrpController/pkg/info"
	_type "LoCyanFrpController/pkg/type"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// NewWebSocket 初始化WebSocket客户端
func NewWebSocket() *WsClient {
	log.Print("Start to Connect WebSocket...")
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
		fmt.Printf("Failed to send message, err: %v", err)
	}

	return nil
}

// ReadMsg 从服务器读取消息
func (w *WsClient) ReadMsg() {
	defer func() {
		err := w.conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	for {
		_, msg, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
		log.Printf("Received message from server: %v", msg)
		var msgJson WsResponse
		err = json.Unmarshal(msg, &msgJson)
		if err != nil {
			log.Printf("Cant unmarshal json message, err: %v", err)
		}
		if msgJson.Status != 200 {
			log.Printf("Error Message from server: %v", msgJson.Message)
		}
	}
}

func (w *WsClient) sendNodeInfoToServer(serverInfo _type.FrpsServerInfoResponse) {
	// nodeInfo
	err := w.SendMsg("upload-node-stats", info.GetNodeInfo(serverInfo))
	if err != nil {
		log.Fatalf("Send node info to server failed! err: %s", err)
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

func main() {
	ws := NewWebSocket()
	err := ws.ConnectWsServer()
	if err != nil {
		log.Fatalf("Can't connect to WebSocket server, err: %v", err)
	} else {
		log.Printf("Connect to WebSocket server successfully!")
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				log.Fatalf("Can't close WebSocket Connection, err: %v", err)
			}
		}(ws.conn)
		go ws.ReadMsg()
		ticker := time.NewTicker(time.Second * ws.config.SendDuration)
		defer ticker.Stop()

		serverInfo := server.GetServerInfo()
		ws.sendNodeInfoToServer(serverInfo)

		for range ticker.C {
			serverInfo := server.GetServerInfo()
			ws.sendNodeInfoToServer(serverInfo)
		}
	}
}
