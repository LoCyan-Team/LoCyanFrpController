//go:build !linux

package inject

import (
	"context"
	"lcf-controller/logger"
	"lcf-controller/pkg/config"
)

// RunOpenGFW 运行 OpenGFW 引擎
func RunOpenGFW(_ context.Context, _ config.OpenGFWConfig) {
	logger.Logger.Fatal("OpenGFW engine is not supported on this platform")
}
