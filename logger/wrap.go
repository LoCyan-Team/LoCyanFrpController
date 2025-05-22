/**
Wrapped zap logger methods
*/

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Sugar() *zap.SugaredLogger {
	return zapLogger.Sugar()
}
func Named(s string) *zap.Logger {
	return zapLogger.Named(s)
}
func WithOptions(opts ...zap.Option) *zap.Logger {
	return zapLogger.WithOptions(opts...)
}
func With(fields ...zap.Field) *zap.Logger {
	return zapLogger.With(fields...)
}
func WithLazy(fields ...zap.Field) *zap.Logger {
	return zapLogger.WithLazy(fields...)
}
func Level() zapcore.Level {
	return zapLogger.Level()
}
func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return zapLogger.Check(lvl, msg)
}
func Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	zapLogger.Log(lvl, msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}
func Info(msg string, fields ...zap.Field) {
	zapLogger.Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	zapLogger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	zapLogger.Error(msg, fields...)
}
func DPanic(msg string, fields ...zap.Field) {
	zapLogger.DPanic(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	zapLogger.Panic(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	zapLogger.Fatal(msg, fields...)
}
func Sync() error {
	return zapLogger.Sync()
}
func Core() zapcore.Core {
	return zapLogger.Core()
}
func Name() string {
	return zapLogger.Name()
}
