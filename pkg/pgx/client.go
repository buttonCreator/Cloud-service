package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

const healthCheckTimeoutSeconds = 2

type Client struct {
	cfg  config
	conn *pgxpool.Pool
}

func New(connString string, options ...Option) *Client {
	cfg := defaultConfig
	for _, option := range options {
		option(&cfg)
	}

	client := &Client{
		cfg: cfg,
	}
	client.cfg.DSN = connString

	return client
}

func (c *Client) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeoutSeconds*time.Second)
	defer cancel()

	if err := c.conn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping pgx: %w", err)
	}

	return nil
}

func (c *Client) Start(ctx context.Context) error {
	return c.connect(ctx)
}

func (c *Client) Shutdown() error {
	c.conn.Close()

	return nil
}

func (c *Client) Conn() *pgxpool.Pool {
	return c.conn
}

func (c *Client) connect(ctx context.Context) error {
	cfg, err := pgxpool.ParseConfig(c.cfg.DSN)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	cfg.HealthCheckPeriod = c.cfg.HealthCheckPeriod
	cfg.ConnConfig.ConnectTimeout = c.cfg.ConnectTimeout

	tracers := []pgx.QueryTracer{
		otelpgx.NewTracer(),
	}

	if c.cfg.Logger != nil {
		traceLogger := tracelog.TraceLog{
			Logger:   c.cfg.Logger,
			LogLevel: c.cfg.LogLevel,
		}

		tracers = append(tracers, &traceLogger)
	}

	cfg.ConnConfig.Tracer = &multiQueryTracer{Tracers: tracers}

	maxAttempts := c.cfg.MaxAttempts
	if maxAttempts == 0 {
		maxAttempts = defaultConfig.MaxAttempts
	}

	pool, err := c.connectWithRetries(ctx, maxAttempts, cfg)
	if err != nil {
		return fmt.Errorf("connect master: %w", err)
	}
	c.conn = pool

	return nil
}

func (c *Client) connectWithRetries(
	ctx context.Context,
	connectAttempts uint,
	cfg *pgxpool.Config,
) (*pgxpool.Pool, error) {
	var err error
	for i := uint(0); i < connectAttempts; i++ {
		if i != 0 {
			time.Sleep(c.cfg.RetryDelay)
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to connect pgx: %w", ctx.Err())
		default:
		}

		var pool *pgxpool.Pool
		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			return pool, nil
		}

		c.cfg.Logger.With("attempt", i+1).WithError(err).Warn("Connect to PostgresSQL")
	}

	return nil, fmt.Errorf("connect to PostgresSQL repository with %d attempts: %w", connectAttempts, err)
}
