package tracing

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const maxStatementLength = 500

type TracedDBTX struct {
	db DBTX
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func NewTracedDBTX(db DBTX) *TracedDBTX {
	return &TracedDBTX{db: db}
}

func truncateStatement(stmt string) string {
	if len(stmt) > maxStatementLength {
		return stmt[:maxStatementLength] + "..."
	}
	return stmt
}

func (t *TracedDBTX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx, span := Tracer().Start(ctx, "db.ExecContext",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", truncateStatement(query)),
		),
	)
	defer span.End()

	result, err := t.db.ExecContext(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}

func (t *TracedDBTX) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	ctx, span := Tracer().Start(ctx, "db.PrepareContext",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", truncateStatement(query)),
		),
	)
	defer span.End()

	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return stmt, err
}

func (t *TracedDBTX) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx, span := Tracer().Start(ctx, "db.QueryContext",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", truncateStatement(query)),
		),
	)
	defer span.End()

	rows, err := t.db.QueryContext(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return rows, err
}

func (t *TracedDBTX) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	_, span := Tracer().Start(ctx, "db.QueryRowContext",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", truncateStatement(query)),
		),
	)
	defer span.End()

	return t.db.QueryRowContext(ctx, query, args...)
}
