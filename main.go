package main

import (
	"context"
	"lcf-controller/inject"
	"lcf-controller/logger"
	"lcf-controller/pkg/config"
	"lcf-controller/server"
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
		err := server.SendTunnelTrafficToServer(cfg)
		if err != nil {
			logger.Error("Can't send proxy traffic to server", zap.Error(err))
		}
		ticker := time.NewTicker(cfg.ControllerConfig.SendDuration)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				logger.Info("shutting down...")
				return
			case <-ticker.C:
				err := server.SendTunnelTrafficToServer(cfg)
				if err != nil {
					logger.Error("Can't send proxy traffic to server", zap.Error(err))
				}
			}
		}
	}
}
