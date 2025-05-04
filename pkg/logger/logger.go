package logger

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ ILogger = (ILogger)(nil)

// ILogger interface
type ILogger interface {
	Ctx(ctx context.Context) ILogger
	With(args ...any) ILogger
	WithError(e error) ILogger
	Sync() error
	Debug(msg string, keyAndValues ...any)
	Info(msg string, keyAndValues ...any)
	Warn(msg string, keyAndValues ...any)
	Error(msg string, keyAndValues ...any)
	Fatal(msg string, keyAndValues ...any)
	Panic(msg string, keyAndValues ...any)
	Log(ctx context.Context, _ tracelog.LogLevel, msg string, data map[string]any)
}

// Logger struct
type Logger struct {
	config config

	SugaredLogger *otelzap.SugaredLogger
	ctx           context.Context
}

const zapCallerSkip = 2

func (l *Logger) Log(ctx context.Context, _ tracelog.LogLevel, msg string, data map[string]any) {
	l.SugaredLogger.InfowContext(ctx, msg, data)
}

// New init zap logger
func New(options ...Option) (*Logger, error) {
	l := Logger{}
	for _, option := range options {
		option(&l.config)
	}

	var zapConfig zap.Config
	zapConfig = zap.NewDevelopmentConfig()
	zapConfig.Sampling = nil
	zapConfig.Level = zap.NewAtomicLevelAt(l.config.level)

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("init logger %w", err)
	}

	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(zapCallerSkip))

	otelLogger := otelzap.New(zapLogger, otelzap.WithMinLevel(zapcore.InfoLevel))
	l.SugaredLogger = otelLogger.Sugar()

	return &l, nil
}

// Ctx function
func (l *Logger) Ctx(ctx context.Context) ILogger {
	fieldsValuesSlice := contextFields(ctx)

	ll := l.With(fieldsValuesSlice...).(*Logger)
	ll.ctx = ctx

	return ll
}

// With function
func (l *Logger) With(args ...any) ILogger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(args...),
		ctx:           l.ctx,
	}
}

// WithError function
func (l *Logger) WithError(e error) ILogger {
	return l.With(zap.Error(e))
}

// Sync function
func (l *Logger) Sync() error {
	if err := l.SugaredLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		return fmt.Errorf("sugared logger sync error: %w", err)
	}

	return nil
}

// Debug function
func (l *Logger) Debug(msg string, keyAndValues ...any) {
	l.SugaredLogger.DebugwContext(l.ctx, msg, keyAndValues...)
}

// Info function
func (l *Logger) Info(msg string, keyAndValues ...any) {
	l.SugaredLogger.InfowContext(l.ctx, msg, keyAndValues...)
}

// Warn function
func (l *Logger) Warn(msg string, keyAndValues ...any) {
	l.SugaredLogger.WarnwContext(l.ctx, msg, keyAndValues...)
}

// Error function
func (l *Logger) Error(msg string, keyAndValues ...any) {
	l.SugaredLogger.ErrorwContext(l.ctx, msg, keyAndValues...)
}

// Fatal function
func (l *Logger) Fatal(msg string, keyAndValues ...any) {
	l.SugaredLogger.FatalwContext(l.ctx, msg, keyAndValues...)
}

// Panic function
func (l *Logger) Panic(msg string, keyAndValues ...any) {
	l.SugaredLogger.PanicwContext(l.ctx, msg, keyAndValues...)
}
