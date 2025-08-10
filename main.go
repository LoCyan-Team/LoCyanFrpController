package main

import (
	"context"
	"lcf-controller/inject"
	"lcf-controller/logger"
	"lcf-controller/pkg/config"
	websocket2 "lcf-controller/pkg/websocket"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func createContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Graceful shutdown
		shutdownChan := make(chan os.Signal, 1)
		signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
		<-shutdownChan
		logger.Info("shutting down gracefully...")
		cancel()
	}()
	return ctx, cancel
}

func main() {
	if runtime.GOOS != "windows" && os.Getuid() != 0 {
		logger.Fatal("please run as root user")
		return
	}

	cfg := config.ReadCfg()

	ctx, _ := createContext()
	if cfg.OpenGFWConfig.Enable {
		go inject.RunOpenGFW(ctx, cfg.OpenGFWConfig)
	}
	if cfg.MonitorConfig.Enable {
		go inject.RunAkileMonitor(ctx, cfg.MonitorConfig)
	}

	if cfg.ControllerConfig.Enable {
		ws := websocket2.NewWebSocket(ctx)
		logger.Info("connecting to WebSocket endpoint...")
		err := ws.ConnectWsServer()
		if err != nil {
			logger.Fatal(
				"can't connect to WebSocket server",
				zap.Error(err),
			)
		} else {
			logger.Info("connect to WebSocket server successfully")

			// 链接成功后就可以开始接收消息了
			// 异步处理
			go func() {
				err := ws.ReadMsg()
				if err != nil {
					logger.Fatal("Cannot read message", zap.Error(err))
				}
			}()
			defer ws.Disconnect()

			// 订阅消息
			err := ws.Subscribe("/traffic")
			if err != nil {
				logger.Fatal("Cannot subscribe traffic", zap.Error(err))
			}

			// 先发一个
			ws.SendProxyStatsToServer(cfg)
			ticker := time.NewTicker(cfg.ControllerConfig.SendDuration)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					logger.Info("shutting down...")
					return
				case <-ticker.C:
					ws.SendProxyStatsToServer(cfg)
				}
			}
		}
	}
}
