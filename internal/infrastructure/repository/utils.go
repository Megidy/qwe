package repository

import (
	"errors"

	"github.com/Megidy/cats/internal/domain/constants"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == constants.PgErrForeignKeyViolation
	}
	return false
}

func isUniqueKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == constants.PgErrUniqueKeyViolation
	}
	return false
}
func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
