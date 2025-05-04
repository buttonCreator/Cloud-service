package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type multiQueryTracer struct {
	Tracers []pgx.QueryTracer
}

func (m *multiQueryTracer) TraceQueryStart(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	for _, trace := range m.Tracers {
		ctx = trace.TraceQueryStart(ctx, conn, data)
	}

	return ctx
}

func (m *multiQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, trace := range m.Tracers {
		trace.TraceQueryEnd(ctx, conn, data)
	}
}
