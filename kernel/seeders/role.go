package seeders

import (
	"context"
	"labflow/models"
	"labflow/sqlc"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedRole(ctx context.Context, pool *pgxpool.Pool) {
	q := sqlc.New(pool)
	roles := []models.Role{
		{ID: 1, RoleName: func(s string) *string { return &s }("admin")},
		{ID: 2, RoleName: func(s string) *string { return &s }("customer")},
		{ID: 3, RoleName: func(s string) *string { return &s }("qa")},
		{ID: 4, RoleName: func(s string) *string { return &s }("micro-analyst")},
		{ID: 5, RoleName: func(s string) *string { return &s }("chemical-analyst")},
		{ID: 6, RoleName: func(s string) *string { return &s }("micro-qc")},
		{ID: 7, RoleName: func(s string) *string { return &s }("chemical-qc")},
		{ID: 8, RoleName: func(s string) *string { return &s }("lab-head")},
		{ID: 9, RoleName: func(s string) *string { return &s }("mr")},
	}
	for _, r := range roles {
		if err := q.SeedRole(ctx, sqlc.SeedRoleParams{RoleID: r.ID, RoleName: r.RoleName}); err != nil {
			log.Fatalf("Failed to seed role data: %v", err)
		}
	}
}
