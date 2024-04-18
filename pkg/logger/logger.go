package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Init(loggerLevel string) {
	// Логгирование
	atomicLevel := zap.NewAtomicLevel()
	switch loggerLevel {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "fatal":
		atomicLevel.SetLevel(zap.FatalLevel)
	case "info":
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	config := zap.Config{
		Level:       atomicLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "log.json"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ = config.Build()
	defer logger.Sync()
}

func Error(text string, err error) {
	logger.Error(text, zap.Error(err))
}

func Warn(text string, err error) {
	logger.Warn(text, zap.Error(err))
}
