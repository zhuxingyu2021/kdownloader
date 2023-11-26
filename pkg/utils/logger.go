package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func initLogger() {
	_, use_stdout := os.LookupEnv("USE_STDOUT")
	var outPutPaths []string
	if use_stdout {
		outPutPaths = append(outPutPaths, "stdout")
	}
	// 配置 logger
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "json",
		OutputPaths:      append(outPutPaths, "/tmp/kdl/logs"),
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
