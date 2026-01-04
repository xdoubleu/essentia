package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xdoubleu/essentia/pkg/database"
)

// SpanDB is used to wrap database actions in [sentry.Span]s.
type SpanDB struct {
	DB     DB
	dbName string
}

// NewSpanDB creates a new [*SpanDB].
func NewSpanDB(db DB) *SpanDB {
	return &SpanDB{
		DB:     db,
		dbName: "postgresql",
	}
}

// Exec is used to wrap Exec in a [sentry.Span].
func (db *SpanDB) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (pgconn.CommandTag, error) {
	return database.WrapWithSpan(ctx, db.dbName, db.DB.Exec, sql, arguments...)
}

// Query is used to wrap Query in a [sentry.Span].
func (db *SpanDB) Query(
	ctx context.Context,
	sql string,
	optionsAndArgs ...any,
) (pgx.Rows, error) {
	return database.WrapWithSpan(ctx, db.dbName, db.DB.Query, sql, optionsAndArgs...)
}

// QueryRow is used to wrap QueryRow in a [sentry.Span].
func (db *SpanDB) QueryRow(
	ctx context.Context,
	sql string,
	optionsAndArgs ...any,
) pgx.Row {
	return database.WrapWithSpanNoError(
		ctx,
		db.dbName,
		db.DB.QueryRow,
		sql,
		optionsAndArgs...)
}

// SendBatch is used to wrap SendBatch in a [sentry.Span].
func (db *SpanDB) SendBatch(
	ctx context.Context,
	b *pgx.Batch,
) pgx.BatchResults {
	sql := ""
	for i, query := range b.QueuedQueries {
		sql += fmt.Sprintf("query %d: %s\n", i, query.SQL)
	}

	span := database.StartSpan(ctx, db.dbName, sql)
	defer span.Finish()

	return db.DB.SendBatch(ctx, b)
}

// Begin doesn't wrap Begin in a [sentry.Span] as
// this makes little sense for starting a transaction.
func (db *SpanDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.DB.Begin(ctx)
}

// BeginTx doesn't wrap BeginTx in a [sentry.Span] as
// this makes little sense for starting a transaction.
func (db *SpanDB) BeginTx(
	ctx context.Context,
	txOptions pgx.TxOptions,
) (pgx.Tx, error) {
	return db.DB.BeginTx(ctx, txOptions)
}

// Ping doesn't wrap Ping in a [sentry.Span] as
// this makes little sense for pinging the db.
func (db *SpanDB) Ping(ctx context.Context) error {
	return db.DB.Ping(ctx)
}
