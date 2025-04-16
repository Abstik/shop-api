package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// 初始化zap日志

// 初始化zap日志
func InitLogger() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // 时间格式
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 使用 zap 内置颜色支持
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 控制台格式
		zapcore.Lock(os.Stdout),                  // 输出到标准输出
		zap.NewAtomicLevelAt(zapcore.DebugLevel), // 允许 DEBUG 及以上日志
	)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}
