package seeders

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunSeeders(db *pgxpool.Pool) {
	ctx := context.Background()
	SeedAccess(ctx, db)
	SeedRole(ctx, db)
	SeedRoleAccess(ctx, db)
	SeedUser(ctx, db)
	SeedCategory(ctx, db)
	SeedTestList(ctx, db)
}
