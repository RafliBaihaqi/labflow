package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"labflow/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{pool: pool, q: sqlc.New(pool)}
}

func (r *AuthRepo) CreateUser(ctx context.Context, username, password, email string, phoneNumber string, name string, company, address string, roleID int) (*models.User, error) {
	var user models.User
	err := inTx(ctx, r.pool, func(q *sqlc.Queries) error {
		row, err := q.CreateUser(ctx, sqlc.CreateUserParams{
			Xid:         uuid.New(),
			Username:    username,
			Password:    password,
			Email:       email,
			PhoneNumber: &phoneNumber,
			Name:        name,
			Company:     &company,
			Address:     &address,
		})
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", wrapDBError(err))
		}
		user = toUser(row)

		if roleID > 0 {
			if err := q.AddUserRole(ctx, sqlc.AddUserRoleParams{UserID: user.ID, RoleID: roleID}); err != nil {
				return fmt.Errorf("failed to map user role: %w", wrapDBError(err))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, wrapDBError(err)
	}
	user := toUser(row)
	return &user, nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, wrapDBError(err)
	}
	user := toUser(row)
	return &user, nil
}

func (r *AuthRepo) UpdateUserVerifiedStatus(ctx context.Context, email string) (*models.User, error) {
	row, err := r.q.SetUserVerified(ctx, email)
	if err != nil {
		return nil, wrapDBError(err)
	}
	user := toUser(row)
	return &user, nil
}

// GetUserByID returns the user with their roles and each role's accesses.
func (r *AuthRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	row, err := r.q.GetUserByID(ctx, userID)
	if err != nil {
		return nil, wrapDBError(err)
	}
	user := toUser(row)

	perms, err := r.q.ListUserPermissions(ctx, userID)
	if err != nil {
		return nil, wrapDBError(err)
	}
	for _, p := range perms {
		if n := len(user.Roles); n == 0 || user.Roles[n-1].ID != p.RoleID {
			user.Roles = append(user.Roles, models.Role{ID: p.RoleID, RoleName: p.RoleName})
		}
		if p.AccessID != nil {
			role := &user.Roles[len(user.Roles)-1]
			role.Accesses = append(role.Accesses, models.Access{ID: *p.AccessID, AccessName: p.AccessName})
		}
	}
	return &user, nil
}

func (r *AuthRepo) UpdateUserPassword(ctx context.Context, email string, hashedPassword string) (*models.User, error) {
	row, err := r.q.SetUserPassword(ctx, sqlc.SetUserPasswordParams{Email: email, Password: hashedPassword})
	if err != nil {
		return nil, wrapDBError(err)
	}
	user := toUser(row)
	return &user, nil
}
