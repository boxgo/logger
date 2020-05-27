package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newAtomicLevelAt(lvStr string) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(levelString2Level(lvStr))
}

func levelString2Level(lvStr string) zapcore.Level {
	lv := zapcore.InfoLevel

	switch lvStr {
	case "debug":
		lv = zapcore.DebugLevel
	case "info":
		lv = zapcore.InfoLevel
	case "warn":
		lv = zapcore.WarnLevel
	case "error":
		lv = zapcore.ErrorLevel
	case "dpanic":
		lv = zapcore.DPanicLevel
	case "panic":
		lv = zapcore.PanicLevel
	case "fatal":
		lv = zapcore.FatalLevel
	}

	return lv
}
