package main

import (
	"LoCyanFrpController/pkg/config"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

// NewWebSocket 初始化WebSocket客户端
func NewWebSocket() *WsClient {
	fmt.Print("Start to Connect WebSocket...")
	ws := new(WsClient)
	cfgInfo := config.ReadCfg()
	ws.addr = cfgInfo.Addr
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
func (w *WsClient) SendMsg(msg string) (err error) {
	if err := w.conn.WriteMessage(websocket.BinaryMessage, []byte(msg)); err != nil {
		log.Printf("Send message successfully! Content: %v", msg)
	} else {
		log.Printf("Failed to send message: %v", err)
	}
	return err
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
	}
}

// WsClient WebSocket客户端结构
type WsClient struct {
	addr string
	conn *websocket.Conn
}

func main() {
	ws := NewWebSocket()
	err := ws.ConnectWsServer()
	if err != nil {
		log.Fatalf("Can't connect to WebSocket server, err: %v", err)
	} else {
		log.Printf("Connect to WebSocket server successfully!")
		go ws.ReadMsg()

		err = ws.SendMsg("Hello, world!")
		if err != nil {
			return
		}
	}
}
