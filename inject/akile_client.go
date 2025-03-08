package inject

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"github.com/henrylee2cn/goutil/calendar/cron"
	"go.uber.org/zap"
	"lcf-controller/inject/akile_monitor_client"
	"lcf-controller/inject/akile_monitor_client/model"
	"lcf-controller/logger"
	"lcf-controller/pkg/config"
	"net/http"
	"time"
)

var conn *websocket.Conn

// 封装连接函数
func connectEndpoint(cfg config.MonitorConfig) (ws *websocket.Conn, err error) {
	logger.Logger.Info("connecting to status endpoint...")

	headers := http.Header{}
	headers.Set("User-Agent", "LoCyanFrp/1.0 (Controller; Status Report)")
	c, _, err := websocket.DefaultDialer.Dial(cfg.Addr, headers)
	return c, err
}

func RunAkileMonitor(ctx context.Context, cfg config.MonitorConfig) {
	go func() {
		c := cron.New()
		if err := c.AddFunc("* * * * * *", func() {
			akile_monitor_client.TrackNetworkSpeed()
		}); err != nil {
			logger.Logger.Fatal("failed to run monitor cronjob", zap.Error(err))
		}
		c.Start()
	}()

	flag.Parse()

	connect := func() {
		for {
			wsc, err := connectEndpoint(cfg)
			if err != nil {
				logger.Logger.Error("error dial status endpoint", zap.Error(err))
				time.Sleep(5 * time.Second)
				continue
			}
			conn = wsc

			_ = conn.WriteMessage(websocket.TextMessage, []byte(cfg.AuthSecret))
			_, message, err := conn.ReadMessage()
			if err != nil {
				logger.Logger.Error("error while status endpoint authentication", zap.Error(err))
				conn.Close()
				time.Sleep(5 * time.Second)
				continue
			}

			if string(message) == "auth success" {
				logger.Logger.Info("connect to status endpoint successfully")
				break
			} else {
				logger.Logger.Error("status endpoint authentication failed, please check your configuration", zap.Error(err))
				conn.Close()
				time.Sleep(5 * time.Second)
				continue
			}
		}
	}

	// 初始化连接
	connect()
	defer conn.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			var D struct {
				Host      *model.Host
				State     *model.HostState
				TimeStamp int64
			}
			D.Host = akile_monitor_client.GetHost()
			D.State = akile_monitor_client.GetState()
			D.TimeStamp = t.Unix()

			// gzip压缩json
			dataBytes, err := json.Marshal(D)
			if err != nil {
				logger.Logger.Error("json.Marshal error", zap.Error(err))
				continue
			}

			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			if _, err := gz.Write(dataBytes); err != nil {
				logger.Logger.Error("gzip write error", zap.Error(err))
				continue
			}

			if err := gz.Close(); err != nil {
				logger.Logger.Error("gzip close error", zap.Error(err))
				continue
			}

			// 发送数据
			err = conn.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				logger.Logger.Error("reporting server status to endpoint error", zap.Error(err))
				conn.Close()
				connect()
			}
		}
	}
}
