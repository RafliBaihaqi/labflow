package repo

import (
	"errors"
	"labflow/apperr"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func wrapDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return apperr.ErrDuplicateKey
	}
	return err
}
