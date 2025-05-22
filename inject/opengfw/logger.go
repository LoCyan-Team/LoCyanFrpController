package opengfw

import (
	"fmt"
	"github.com/apernet/OpenGFW/ruleset"
	"go.uber.org/zap"
	"lcf-controller/logger"
)

type EngineLogger struct{}

func (l *EngineLogger) WorkerStart(id int) {
	logger.Debug("worker started", zap.Int("id", id))
}

func (l *EngineLogger) WorkerStop(id int) {
	logger.Debug("worker stopped", zap.Int("id", id))
}

func (l *EngineLogger) TCPStreamNew(workerID int, info ruleset.StreamInfo) {
	logger.Debug("new TCP stream",
		zap.Int("workerID", workerID),
		zap.Int64("id", info.ID),
		zap.String("src", info.SrcString()),
		zap.String("dst", info.DstString()))
}

func (l *EngineLogger) TCPStreamPropUpdate(info ruleset.StreamInfo, close bool) {
	logger.Debug("TCP stream property update",
		zap.Int64("id", info.ID),
		zap.String("src", info.SrcString()),
		zap.String("dst", info.DstString()),
		zap.Any("props", info.Props),
		zap.Bool("close", close))
}

func (l *EngineLogger) TCPStreamAction(info ruleset.StreamInfo, action ruleset.Action, noMatch bool) {
	if noMatch {
		logger.Debug("TCP stream no match",
			zap.Int64("id", info.ID),
			zap.String("src", info.SrcString()),
			zap.String("dst", info.DstString()),
			zap.String("action", action.String()))
	} else {
		logger.Info("TCP stream action",
			zap.Int64("id", info.ID),
			zap.String("src", info.SrcString()),
			zap.String("dst", info.DstString()),
			zap.String("action", action.String()))
	}
}

func (l *EngineLogger) TCPFlush(workerID, flushed, closed int) {
	logger.Debug("TCP flush",
		zap.Int("workerID", workerID),
		zap.Int("flushed", flushed),
		zap.Int("closed", closed))
}

func (l *EngineLogger) UDPStreamNew(workerID int, info ruleset.StreamInfo) {
	logger.Debug("new UDP stream",
		zap.Int("workerID", workerID),
		zap.Int64("id", info.ID),
		zap.String("src", info.SrcString()),
		zap.String("dst", info.DstString()))
}

func (l *EngineLogger) UDPStreamPropUpdate(info ruleset.StreamInfo, close bool) {
	logger.Debug("UDP stream property update",
		zap.Int64("id", info.ID),
		zap.String("src", info.SrcString()),
		zap.String("dst", info.DstString()),
		zap.Any("props", info.Props),
		zap.Bool("close", close))
}

func (l *EngineLogger) UDPStreamAction(info ruleset.StreamInfo, action ruleset.Action, noMatch bool) {
	if noMatch {
		logger.Debug("UDP stream no match",
			zap.Int64("id", info.ID),
			zap.String("src", info.SrcString()),
			zap.String("dst", info.DstString()),
			zap.String("action", action.String()))
	} else {
		logger.Info("UDP stream action",
			zap.Int64("id", info.ID),
			zap.String("src", info.SrcString()),
			zap.String("dst", info.DstString()),
			zap.String("action", action.String()))
	}
}

func (l *EngineLogger) ModifyError(info ruleset.StreamInfo, err error) {
	logger.Error("modify error",
		zap.Int64("id", info.ID),
		zap.String("src", info.SrcString()),
		zap.String("dst", info.DstString()),
		zap.Error(err))
}

func (l *EngineLogger) AnalyzerDebugf(streamID int64, name string, format string, args ...interface{}) {
	logger.Debug("analyzer debug message",
		zap.Int64("id", streamID),
		zap.String("name", name),
		zap.String("msg", fmt.Sprintf(format, args...)))
}

func (l *EngineLogger) AnalyzerInfof(streamID int64, name string, format string, args ...interface{}) {
	logger.Info("analyzer info message",
		zap.Int64("id", streamID),
		zap.String("name", name),
		zap.String("msg", fmt.Sprintf(format, args...)))
}

func (l *EngineLogger) AnalyzerErrorf(streamID int64, name string, format string, args ...interface{}) {
	logger.Error("analyzer error message",
		zap.Int64("id", streamID),
		zap.String("name", name),
		zap.String("msg", fmt.Sprintf(format, args...)))
}

type RulesetLogger struct{}

func (l *RulesetLogger) Log(info ruleset.StreamInfo, name string) {
	logger.Info("ruleset log",
		zap.String("name", name),
		zap.Int64("id", info.ID),
		zap.String("src", info.SrcString()),
		zap.String("dst", info.DstString()),
		zap.Any("props", info.Props))
}

func (l *RulesetLogger) MatchError(info ruleset.StreamInfo, name string, err error) {
	logger.Error("ruleset match error",
		zap.String("name", name),
		zap.Int64("id", info.ID),
		zap.String("src", info.SrcString()),
		zap.String("dst", info.DstString()),
		zap.Error(err))
}
