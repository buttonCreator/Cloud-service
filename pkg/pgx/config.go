package pgx

import (
	"time"

	"github.com/jackc/pgx/v5/tracelog"

	"Cloud/pkg/logger"
)

type config struct {
	DSN                      string
	ConnectTimeout           time.Duration
	ConnectionRetriesTimeout time.Duration
	HealthCheckPeriod        time.Duration
	MaxAttempts              uint
	RetryDelay               time.Duration

	Logger   logger.ILogger
	LogLevel tracelog.LogLevel
}

const (
	retryDelay     = 2
	defaultTimeout = 5
	maxAttempts    = 5
)

var defaultConfig = config{
	ConnectTimeout:           defaultTimeout * time.Second,
	ConnectionRetriesTimeout: defaultTimeout * time.Second,
	HealthCheckPeriod:        defaultTimeout * time.Second,

	MaxAttempts: maxAttempts,
	RetryDelay:  retryDelay,
}

type Option func(c *config)

func WithLogger(l logger.ILogger) Option {
	return func(c *config) {
		c.Logger = l
	}
}

func WithLoggerLevel(logLevel tracelog.LogLevel) Option {
	return func(c *config) {
		c.LogLevel = logLevel
	}
}
