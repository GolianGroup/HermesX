package logging

import "go.uber.org/zap"

// ZapLogger is a wrapper around zap.Logger
type ZapLogger struct {
	*zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{Logger: logger}
}

func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.Logger.Info(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.Logger.Error(msg, fields...)
}

func (z *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	z.Logger.Fatal(msg, fields...)
}

// ...
