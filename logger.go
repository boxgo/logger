package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type (
	// Logger logger option
	Logger struct {
		Level          string `config:"level" help:"Levels: debug,info,warn,error,dpanic,panic,fatal"`
		Encoding       string `config:"encoding" help:"PS: console or json"`
		TraceUID       string `config:"traceUid" help:"Name as trace uid in context"`
		TraceRequestID string `config:"traceRequestId" help:"Name as trace requestId in context"`
		CallerSkip     int    `config:"callerSkip" help:"AddCallerSkip increases the number of callers skipped by caller annotation"`
		*zap.SugaredLogger
	}
)

var (
	// Default the default logger
	Default = &Logger{}
	global  = Default
	once    = sync.Once{}
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

	log, err := simpleConfig(logger.Level, logger.Encoding).Build(
		zap.AddCallerSkip(logger.CallerSkip),
	)
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	logger.SugaredLogger = log.Sugar()
}

// ConfigDidLoad did load
func (logger *Logger) ConfigDidLoad(ctx context.Context) {
	logger.ConfigWillLoad(ctx)

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

// Trace logger with requestId and uid
func (logger *Logger) Trace(ctx context.Context) *zap.SugaredLogger {
	return trace(ctx, logger)
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
	uid := ctx.Value(logger.TraceUID)
	requestID := ctx.Value(logger.TraceRequestID)

	// return logger.SugaredLogger.Named(fmt.Sprintf("[%s][%s]", uid, requestID))
	return logger.SugaredLogger.With("uid", uid).With("requestId", requestID)
}
