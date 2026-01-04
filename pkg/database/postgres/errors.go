package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xdoubleu/essentia/pkg/database"
)

// PgxErrorToHTTPError converts a database error
// from [github.com/jackc/pgx] to an appropriate HTTP error.
func PgxErrorToHTTPError(err error) error {
	var pgxError *pgconn.PgError
	errors.As(err, &pgxError)

	switch {
	case pgxError == nil:
		if errors.Is(err, pgx.ErrNoRows) {
			return database.ErrResourceNotFound
		}
		return err
	case pgxError.Code == pgerrcode.ForeignKeyViolation:
		return database.ErrResourceNotFound
	case pgxError.Code == pgerrcode.UniqueViolation:
		return database.ErrResourceConflict
	default:
		return err
	}
}
