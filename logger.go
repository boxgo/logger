package logger

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger logger option
	Logger struct {
		Level          string   `config:"level" help:"Levels: debug,info,warn,error,dpanic,panic,fatal"`
		Encoding       string   `config:"encoding" help:"PS: console or json"`
		TraceUID       string   `config:"traceUid" help:"Name as trace uid in context"`
		TraceRequestID string   `config:"traceRequestId" help:"Name as trace requestId in context"`
		TraceSpanID    string   `config:"traceSpanId" help:"Name as trace spanId in context"`
		TraceBizID     string   `config:"traceBizId" help:"Name as trace spanId in context"`
		CallerKey      string   `config:"callerKey" help:"caller key in log"`
		CallerSkip     int      `config:"callerSkip" help:"AddCallerSkip increases the number of callers skipped by caller annotation"`
		FilterSpecs    []string `config:"filterSpecs" help:"filter rules, split by ==>. eg: old ==> new"`
		*zap.SugaredLogger
	}
)

var (
	// Default the default logger
	Default = &Logger{}
	global  = Default
	once    = sync.Once{}
	locker  = sync.Mutex{}
)

func init() {
	if Default.SugaredLogger == nil {
		l, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		Default.SugaredLogger = l.Sugar()
	}
}

// Name logger config name
func (logger *Logger) Name() string {
	return "logger"
}

// ConfigWillLoad logger config
func (logger *Logger) ConfigWillLoad(context.Context) {
	if logger.Level == "" {
		logger.Level = "debug"
	}
	if logger.Encoding == "" {
		logger.Encoding = "console"
	}
	if logger.TraceUID == "" {
		logger.TraceUID = "uid"
	}
	if logger.TraceRequestID == "" {
		logger.TraceRequestID = "requestId"
	}
	if logger.TraceSpanID == "" {
		logger.TraceSpanID = "traceSpanId"
	}
	if logger.TraceBizID == "" {
		logger.TraceBizID = "traceBizId"
	}

	logger.apply()
}

// ConfigDidLoad did load
func (logger *Logger) ConfigDidLoad(ctx context.Context) {
	logger.apply()

	once.Do(func() {
		global = &Logger{
			Level:          Default.Level,
			Encoding:       Default.Encoding,
			TraceUID:       Default.TraceUID,
			TraceRequestID: Default.TraceRequestID,
			SugaredLogger:  Default.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
		}
	})
}

func (logger *Logger) apply() {
	locker.Lock()
	defer locker.Unlock()

	cfg := &zap.Config{
		Development: false,
		Level:       newAtomicLevelAt(logger.Level),
		Encoding:    logger.Encoding,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			CallerKey:      logger.CallerKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := cfg.Build(
		zap.AddCallerSkip(logger.CallerSkip),
		zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			var ws *filterWriter
			var enc zapcore.Encoder

			if logger.Encoding == "console" {
				enc = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
			} else {
				enc = zapcore.NewJSONEncoder(cfg.EncoderConfig)
			}

			if logger.Level == "debug" || logger.Level == "info" || logger.Level == "warn" {
				ws = newFilterWriter(os.Stdout, newFilterBySlice(logger.FilterSpecs)...)
			} else {
				ws = newFilterWriter(os.Stderr, newFilterBySlice(logger.FilterSpecs)...)
			}

			return zapcore.NewCore(enc, ws, cfg.Level)
		}),
	)
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	logger.SugaredLogger = log.Sugar()
}

func (logger *Logger) TraceRaw(ctx context.Context) *zap.Logger {
	return trace(ctx, logger).Desugar()
}

// Trace logger with requestId and uid
func (logger *Logger) Trace(ctx context.Context) *zap.SugaredLogger {
	return trace(ctx, logger)
}

func TraceRaw(ctx context.Context) *zap.Logger {
	return trace(ctx, Default).Desugar()
}

func Trace(ctx context.Context) *zap.SugaredLogger {
	return trace(ctx, Default)
}

func DPanic(args ...interface{}) {
	global.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	global.DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	global.DPanicw(msg, keysAndValues...)
}

func Debug(args ...interface{}) {
	global.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	global.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	global.Debugw(msg, keysAndValues...)
}

func Desugar() *zap.Logger {
	return global.Desugar()
}

func Error(args ...interface{}) {
	global.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	global.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	global.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	global.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	global.Fatalw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	global.Info(args...)
}

func Infof(template string, args ...interface{}) {
	global.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	global.Infow(msg, keysAndValues...)
}

func Named(name string) *zap.SugaredLogger {
	return global.Named(name)
}

func Panic(args ...interface{}) {
	global.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	global.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	global.Panicw(msg, keysAndValues...)
}

func Sync() error {
	return global.Sync()
}

func Warn(args ...interface{}) {
	global.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	global.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	global.Warnw(msg, keysAndValues...)
}

func With(args ...interface{}) *zap.SugaredLogger {
	return global.With(args...)
}

func trace(ctx context.Context, logger *Logger) *zap.SugaredLogger {
	var uid, requestID, spanId, bizId string

	if uidStr, ok := ctx.Value(logger.TraceUID).(string); ok {
		uid = uidStr
	}
	if requestIDStr, ok := ctx.Value(logger.TraceRequestID).(string); ok {
		requestID = requestIDStr
	}
	if spanIdStr, ok := ctx.Value(logger.TraceSpanID).(string); ok {
		spanId = spanIdStr
	}
	if bizIdStr, ok := ctx.Value(logger.TraceBizID).(string); ok {
		bizId = bizIdStr
	}

	return logger.SugaredLogger.Named(fmt.Sprintf("[%s][%s][%s][%s]", uid, requestID, spanId, bizId))
}
