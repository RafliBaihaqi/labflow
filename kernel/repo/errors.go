package repo

import (
	"context"
	"errors"
	"fmt"
	"labflow/apperr"
	"labflow/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func wrapDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return apperr.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return apperr.ErrDuplicateKey
	}
	return err
}

// inTx runs fn in a transaction and commits it if fn returns nil.
func inTx(ctx context.Context, pool *pgxpool.Pool, fn func(q *sqlc.Queries) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := fn(sqlc.New(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// pageBounds turns page/limit into LIMIT and OFFSET the way GORM did:
// a negative limit means no limit and a non-positive offset means none.
func pageBounds(page, limit int) (*int, int) {
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		return nil, offset
	}
	return &limit, offset
}

// likePattern returns the ILIKE argument for a search term, or nil for no filter.
func likePattern(search string) *string {
	if search == "" {
		return nil
	}
	pattern := "%" + search + "%"
	return &pattern
}

func deref[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}
