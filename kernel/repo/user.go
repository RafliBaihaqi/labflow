package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"labflow/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	q *sqlc.Queries
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{q: sqlc.New(pool)}
}

func (u *UserRepo) UpdateUserByXid(ctx context.Context, xid uuid.UUID, phoneNumber, company, address, avatar *string, updatedBy int) (*models.User, error) {
	row, err := u.q.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
		Xid:         xid,
		PhoneNumber: phoneNumber,
		Company:     company,
		Address:     address,
		Avatar:      avatar,
		UpdatedBy:   &updatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update reservation: %w", wrapDBError(err))
	}
	updated := toUser(row)
	return &updated, nil
}
