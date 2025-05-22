/**
Wrapped zap logger methods
*/

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Sugar() *zap.SugaredLogger {
	return logger.Sugar()
}
func Named(s string) *zap.Logger {
	return logger.Named(s)
}
func WithOptions(opts ...zap.Option) *zap.Logger {
	return logger.WithOptions(opts...)
}
func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}
func WithLazy(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}
func Level() zapcore.Level {
	return logger.Level()
}
func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return logger.Check(lvl, msg)
}
func Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	logger.Log(lvl, msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
func Sync() error {
	return logger.Sync()
}
func Core() zapcore.Core {
	return logger.Core()
}
func Name() string {
	return logger.Name()
}
