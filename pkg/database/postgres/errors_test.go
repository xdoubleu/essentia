package postgres_test

import (
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/database"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
)

func newPgError(code string) *pgconn.PgError {
	//nolint:exhaustruct //not using other fields
	return &pgconn.PgError{
		Code: code,
	}
}

func TestErrResourceNotFound(t *testing.T) {
	err1 := postgres.PgxErrorToHTTPError(pgx.ErrNoRows)
	err2 := postgres.PgxErrorToHTTPError(
		newPgError(pgerrcode.ForeignKeyViolation),
	)

	assert.ErrorIs(t, err1, database.ErrResourceNotFound)
	assert.ErrorIs(t, err2, database.ErrResourceNotFound)
}

func TestErrResourceConflict(t *testing.T) {
	err := postgres.PgxErrorToHTTPError(
		newPgError(pgerrcode.UniqueViolation),
	)

	assert.ErrorIs(t, err, database.ErrResourceConflict)
}
