package repository

import (
	"github.com/jackc/pgx/v5/tracelog"

	"Cloud/pkg/logger"
	"Cloud/pkg/pgx"
)

type Repository struct {
	*pgx.Client
}

func New(connString string, log logger.ILogger) *Repository {
	return &Repository{
		Client: pgx.New(connString, pgx.WithLogger(log), pgx.WithLoggerLevel(tracelog.LogLevelError)),
	}
}
