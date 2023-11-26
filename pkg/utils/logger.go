package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func initLogger() {
	// 配置 logger
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "json",
		OutputPaths:      []string{"/tmp/kdl/logs", "stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	var err error
	err = os.MkdirAll("/tmp/kdl", 0755)
	if err != nil {
		panic(err)
	}

	// 初始化 logger
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func LoggerSync() {
	err := Logger.Sync()
	if err != nil {
		// Handle error, possibly just:
		// fmt.Fprintf(os.Stderr, "Failed to sync zap logger: %v", err)
		os.Exit(1)
	}
}
