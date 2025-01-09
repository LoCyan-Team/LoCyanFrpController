package websocket

import (
	"LoCyanFrpController/pkg/config"
	"LoCyanFrpController/pkg/info"
	_type "LoCyanFrpController/pkg/type"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

// NewWebSocket 初始化WebSocket客户端
func NewWebSocket() *WsClient {
	log.Print("Start to Connect WebSocket...")
	ws := new(WsClient)
	cfgInfo := config.ReadCfg()
	ws.addr = cfgInfo.Addr
	ws.Config = cfgInfo
	return ws
}

// ConnectWsServer 连接到WebSocket服务器
func (w *WsClient) ConnectWsServer() (err error) {
	conn, _, err := websocket.DefaultDialer.Dial(w.addr, nil)
	if err != nil {
		return err
	}
	w.Conn = conn
	return nil
}

// SendMsg 发送消息到服务器
func (w *WsClient) SendMsg(action string, data map[string]any) (err error) {
	req := new(BasicRequest)
	req.Action = action
	req.Node.Id = w.Config.NodeId
	req.Node.ApiKey = w.Config.NodeApiKey
	req.Data = data
	msg, err := json.Marshal(req)
	err = w.Conn.WriteMessage(websocket.TextMessage, msg)

	if err != nil {
		log.Printf("Failed to send message, err: %v", err)
	}

	return nil
}

// ReadMsg 从服务器读取消息
func (w *WsClient) ReadMsg() {
	defer func() {
		err := w.Conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	for {
		_, msg, err := w.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
		var msgJson WsResponse
		err = json.Unmarshal(msg, &msgJson)
		if err != nil {
			log.Printf("Cant unmarshal json message, err: %v", err)
		}
		if msgJson.Status != 200 {
			log.Printf("Error Message from server: %v", msgJson.Message)
		}
		if msgJson.Status == 200 {
			log.Printf("Received message from server: %v", msgJson.Message)
		}
	}
}

func (w *WsClient) SendNodeStatsToServer(serverInfo _type.FrpsServerInfoResponse) {
	// nodeInfo
	err := w.SendMsg("upload-node-stats", info.GetNodeInfo(serverInfo))
	if err != nil {
		log.Fatalf("Send node info to server failed! err: %s", err)
	}
}

func (w *WsClient) SendProxyStatsToServer() {
	types := []string{"tcp", "udp", "http", "https", "xtcp", "stcp"}
	for _, p := range types {
		proxies := info.GetProxies(p)
		for _, j := range proxies {
			err := w.SendMsg("upload-proxy-stats", j)
			log.Printf("Send proxy info to the server: proxyName: %s, inbound: %v, outbound: %v", j["proxy_name"], j["inbound"], j["outbound"])
			if err != nil {
				log.Fatalf("Send proxy info to server failed! err: %s", err)
			}
		}
	}
}

// WsClient WebSocket客户端结构
type WsClient struct {
	addr   string
	Conn   *websocket.Conn
	Config *config.Config
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
