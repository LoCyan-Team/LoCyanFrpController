package inject

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"github.com/henrylee2cn/goutil/calendar/cron"
	"go.uber.org/zap"
	"lcf-controller/inject/akile_monitor_client"
	"lcf-controller/inject/akile_monitor_client/model"
	"lcf-controller/logger"
	"lcf-controller/pkg/config"
	"os"
	"os/signal"
	"time"
)

func RunAkileMonitor(cfg config.MonitorConfig) {
	go func() {
		c := cron.New()
		c.AddFunc("* * * * * *", func() {
			akile_monitor_client.TrackNetworkSpeed()
		})
		c.Start()
	}()

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	logger.Logger.Info("connecting to status WebSocket endpoint...")

	c, _, err := websocket.DefaultDialer.Dial(cfg.Addr, nil)
	if err != nil {
		logger.Logger.Fatal("error dial status endpoint", zap.Error(err))
	}
	defer c.Close()

	c.WriteMessage(websocket.TextMessage, []byte(cfg.AuthSecret))

	done := make(chan struct{})

	_, message, err := c.ReadMessage()
	if err != nil {
		logger.Logger.Error("invalid auth secret, please check your configuration", zap.Error(err))
		return
	}
	if string(message) == "auth success" {
		logger.Logger.Info("connect to status WebSocket server successfully")
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
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
			//gzip压缩json
			dataBytes, err := json.Marshal(D)
			if err != nil {
				logger.Logger.Error("json.Marshal error", zap.Error(err))
				return
			}

			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			if _, err := gz.Write(dataBytes); err != nil {
				logger.Logger.Error("gzip write error", zap.Error(err))
				return
			}

			if err := gz.Close(); err != nil {
				logger.Logger.Error("gzip close error", zap.Error(err))
				return
			}

			err = c.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				logger.Logger.Error("reporting server status to endpoint error", zap.Error(err))
				RunAkileMonitor(cfg)
				return
			}
		case <-interrupt:
			logger.Logger.Info("closing status endpoint connection...")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Logger.Error("closing server status to endpoint connection error", zap.Error(err))
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
