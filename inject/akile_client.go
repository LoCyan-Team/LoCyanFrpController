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
	logger.Info("connecting to status endpoint...")

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
			logger.Fatal("failed to run monitor cronjob", zap.Error(err))
		}
		c.Start()
	}()

	flag.Parse()

	const maxRetries = 5
connect := func() {
	retries := 0
	for retries < maxRetries {
		wsc, err := connectEndpoint(cfg)
		if err != nil {
			logger.Error("error dial status endpoint", zap.Error(err))
			time.Sleep(5 * time.Second)
			retries++
			continue
		}
		conn = wsc

		_ = conn.WriteMessage(websocket.TextMessage, []byte(cfg.AuthSecret))
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error("error while status endpoint authentication", zap.Error(err))
			conn.Close()
			time.Sleep(5 * time.Second)
			retries++
			continue
		}

		if string(message) == "auth success" {
			logger.Info("connect to status endpoint successfully")
			break
		} else {
			logger.Error("status endpoint authentication failed, please check your configuration", zap.Error(err))
			conn.Close()
			time.Sleep(5 * time.Second)
			retries++
		}
		if retries >= maxRetries {
			logger.Error("max retries reached, giving up connection")
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
				logger.Error("json.Marshal error", zap.Error(err))
				continue
			}

			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			defer gz.Close()
			if _, err := gz.Write(dataBytes); err != nil {
				logger.Error("gzip write error", zap.Error(err))
				continue
			}

			if err := gz.Close(); err != nil {
				logger.Error("gzip close error", zap.Error(err))
				continue
			}

			// 发送数据
			err = conn.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				logger.Error("reporting server status to endpoint error", zap.Error(err))
				conn.Close()
				connect()
			}
		}
	}
}
