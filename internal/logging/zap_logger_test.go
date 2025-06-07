package logging

import (
	"bytes"
	"encoding/json"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLogger(t *testing.T) {
	logger := zap.NewExample()

	type args struct {
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
		want *ZapLogger
	}{
		{
			name: "Valid logger",
			args: args{
				logger: logger,
			},
			want: &ZapLogger{
				Logger: logger,
			},
		},
		{
			name: "Nil logger",
			args: args{
				logger: nil,
			},
			want: &ZapLogger{
				Logger: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewZapLogger(tt.args.logger)
			if (got == nil) != (tt.want == nil) || (got != nil && tt.want != nil && got.Logger != tt.want.Logger) {
				t.Errorf("NewZapLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZapLoggerInfo(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := zap.NewDevelopmentEncoderConfig()

	encoder.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.AddSync(buffer),
		zapcore.InfoLevel,
	)
	logger := zap.New(core)
	zapLogger := NewZapLogger(logger)

	tests := []struct {
		name   string
		logger *ZapLogger
		msg    string
		fields []zap.Field
		want   map[string]interface{}
	}{
		{
			name:   "Simple Info log",
			logger: zapLogger,
			msg:    "Info log message",
			fields: []zap.Field{zap.String("key", "value")},
			want: map[string]interface{}{
				"L":   "INFO",
				"M":   "Info log message",
				"key": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer.Reset()
			tt.logger.Info(tt.msg, tt.fields...)

			got := buffer.String()

			var gotLog map[string]interface{}
			if err := json.Unmarshal([]byte(got), &gotLog); err != nil {
				t.Fatalf("Failed to unmarshal log: %v", err)
			}

			for key, expectedValue := range tt.want {
				if gotLog[key] != expectedValue {
					t.Errorf("Key %q = %v, want %v", key, gotLog[key], expectedValue)
				}
			}

		})
	}
}

func TestZapLoggerError(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := zap.NewDevelopmentEncoderConfig()

	encoder.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.AddSync(buffer),
		zapcore.ErrorLevel,
	)
	logger := zap.New(core)
	zapLogger := NewZapLogger(logger)

	tests := []struct {
		name   string
		logger *ZapLogger
		msg    string
		fields []zap.Field
		want   map[string]interface{}
	}{
		{
			name:   "Simple Error log",
			logger: zapLogger,
			msg:    "Error log message",
			fields: []zap.Field{zap.String("key", "value")},
			want: map[string]interface{}{
				"L":   "ERROR",
				"M":   "Error log message",
				"key": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer.Reset()
			tt.logger.Error(tt.msg, tt.fields...)

			got := buffer.String()

			var gotLog map[string]interface{}
			if err := json.Unmarshal([]byte(got), &gotLog); err != nil {
				t.Fatalf("Failed to unmarshal log: %v", err)
			}

			for key, expectedValue := range tt.want {
				if gotLog[key] != expectedValue {
					t.Errorf("Key %q = %v, want %v", key, gotLog[key], expectedValue)
				}
			}
		})
	}
}
