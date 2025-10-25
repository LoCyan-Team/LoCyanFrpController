package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"lcf-controller/logger"
	"lcf-controller/pkg/config"
	"lcf-controller/pkg/info"
	"lcf-controller/pkg/utils"
	"strings"
	"time"
)

type WsClient struct {
	Addr  string
	WConn *websocket.Conn
	Ctx   context.Context
	Cfg   *config.Config
}

// NewWebSocket 初始化WebSocket客户端
func NewWebSocket(ctx context.Context) *WsClient {
	cfg := config.ReadCfg()
	return &WsClient{
		Addr:  cfg.ControllerConfig.Addr,
		WConn: nil,
		Ctx:   ctx,
		Cfg:   cfg,
	}
}

// ConnectWsServer 连接到WebSocket服务器
func (w *WsClient) ConnectWsServer() (err error) {
	url := w.Addr
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	w.WConn = c

	// 发送 STOMP CONNECT 帧（SockJS 通常与 STOMP 一起使用）
	stompConnect := "CONNECT\naccept-version:1.1,1.0\nheart-beat:10000,10000\n\n\x00"
	err = c.WriteMessage(websocket.TextMessage, []byte(stompConnect))
	if err != nil {
		return err
	}
	// 确保链接
	_, welcomeMsg, err := c.ReadMessage()
	if err != nil || !strings.Contains(string(welcomeMsg), "CONNECTED") {
		return fmt.Errorf("handshake failed: %v", string(welcomeMsg))
	}

	subscribeFrame := "SUBSCRIBE\nid:sub-0\ndestination:/node/greetings\n\n\x00"
	if err := c.WriteMessage(websocket.TextMessage, []byte(subscribeFrame)); err != nil {
		return fmt.Errorf("subscribe failed: %w", err)
	}

	if err := w.SendMsg("/hello", make(map[string]any)); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	_, _, err = c.ReadMessage()
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}
	return nil
}

// SendMsg 发送消息到服务器
func (w *WsClient) SendMsg(destination string, data map[string]any) (err error) {
	// 发送消息
	jsonData, _ := json.Marshal(data)
	jsonDataString := string(jsonData)
	stompSend := fmt.Sprintf("SEND\ndestination:/app%s\ncontent-type:application/json\n\n%s\x00", destination, jsonDataString)
	err = w.WConn.WriteMessage(websocket.TextMessage, []byte(stompSend))
	if err != nil {
		return err
	}
	return nil
}

func (w *WsClient) Subscribe(destination string) (err error) {
	stompSubscribe := fmt.Sprintf("SUBSCRIBE\nid:sub-0\ndestination:/node%s\n\n\x00", destination)
	err = w.WConn.WriteMessage(websocket.TextMessage, []byte(stompSubscribe))
	if err != nil {
		return err
	}
	return nil
}

// ReadMsg 从服务器读取消息
func (w *WsClient) ReadMsg() (err error) {
	// 持续读取消息
	for {
		select {
		case <-w.Ctx.Done():
			logger.Info("stop receive msg...")
			return
		default:
			_, message, err := w.WConn.ReadMessage()
			if err != nil {
				return err
			}
			msgStr := string(message)
			switch {
			case strings.HasPrefix(msgStr, "MESSAGE"):
				dest := utils.ExtractHeader(msgStr, "destination")
				body := utils.ExtractBody(msgStr)
				logger.Info(fmt.Sprintf("Received message from: %s: %s", dest, body))
			default:
				logger.Debug("Received unhandled frame type",
					zap.String("frame", msgStr),
				)
			}
		}
	}
}

func (w *WsClient) SendProxyStatsToServer(cfg *config.Config) {
	types := []string{"tcp", "udp", "http", "https", "xtcp", "stcp"}
	for _, t := range types {
		proxies, err := info.GetProxies(cfg, t)
		if err != nil {
			logger.Error("can't request proxies info", zap.Error(err))
		} else {
			for _, j := range proxies {
				err := w.SendMsg("/traffic", j)
				logger.Info(fmt.Sprintf("Proxy: %s, outBound: %v, inBound: %v", j["proxy_name"], j["in_bound_traffic"], j["out_bound_traffic"]))
				if err != nil {
					logger.Error("send proxy info to server failed!", zap.Error(err))
				}
			}
		}
	}
}

func (w *WsClient) Disconnect() {
	// 发送DISCONNECT帧
	disconnectFrame := "DISCONNECT\nreceipt-id:disconnect-123\n\n\x00"
	if err := w.WConn.WriteMessage(websocket.TextMessage, []byte(disconnectFrame)); err != nil {
		logger.Error("send DISCONNECT failed", zap.Error(err))
	}

	// 等待RECEIPT响应（带错误处理）
	w.WConn.SetReadDeadline(time.Now().Add(2 * time.Second)) // 缩短超时时间
	_, msg, err := w.WConn.ReadMessage()
	if err == nil {
		if strings.HasPrefix(string(msg), "RECEIPT") {
			logger.Info("disconnect receipt confirmed")
		}
	} else if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		logger.Warn("receipt confirmation error", zap.Error(err))
	}

	// 安全关闭连接
	_ = w.WConn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(1*time.Second),
	)
	time.Sleep(100 * time.Millisecond) // 确保关闭指令发送
	w.WConn.Close()
}
