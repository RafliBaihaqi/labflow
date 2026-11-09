package seeders

import (
	"context"
	"labflow/models"
	"labflow/sqlc"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedCategory(ctx context.Context, pool *pgxpool.Pool) {
	q := sqlc.New(pool)
	categories := []models.Category{
		{XID: uuid.New(), Name: "Microbiology", Label: func(s string) *string { return &s }("Microbiology")},
		{XID: uuid.New(), Name: "Chemistry", Label: func(s string) *string { return &s }("Chemistry")},
		{XID: uuid.New(), Name: "Physical", Label: func(s string) *string { return &s }("Physical")},
	}
	for _, c := range categories {
		if err := q.SeedCategory(ctx, sqlc.SeedCategoryParams{Xid: c.XID, Name: c.Name, Label: c.Label}); err != nil {
			log.Fatalf("Failed to seed category data: %v", err)
		}
	}
}
