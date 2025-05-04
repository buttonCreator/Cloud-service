package logger

import (
	"go.uber.org/zap/zapcore"

	commonConfig "Cloud/pkg/config"
)

type config struct {
	env   commonConfig.Env
	level zapcore.Level
}

// Option type
type Option func(*config)

// WithLevel set logging level
func WithLevel(level zapcore.Level) Option {
	return func(config *config) {
		config.level = level
	}
}

// WithEnv set zap env to build logger
func WithEnv(env commonConfig.Env) Option {
	return func(config *config) {
		config.env = env
	}
}
