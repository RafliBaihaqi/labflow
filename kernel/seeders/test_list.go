package seeders

import (
	"context"
	"labflow/models"
	"labflow/sqlc"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedTestList(ctx context.Context, pool *pgxpool.Pool) {
	q := sqlc.New(pool)
	tests := []models.TestList{
		{XID: uuid.New(), CategoryID: 1, Value: "Total Aerobic Microbial Count", Label: "Total Aerobic Microbial Count"},
		{XID: uuid.New(), CategoryID: 1, Value: "Total Yeast and Mold Count", Label: "Total Yeast and Mold Count"},
		{XID: uuid.New(), CategoryID: 1, Value: "Bioburden", Label: "Bioburden"},
		{XID: uuid.New(), CategoryID: 1, Value: "Sterility", Label: "Sterility"},
		{XID: uuid.New(), CategoryID: 1, Value: "Endotoxin", Label: "Endotoxin"},
		{XID: uuid.New(), CategoryID: 1, Value: "E. coli", Label: "E. coli"},
		{XID: uuid.New(), CategoryID: 1, Value: "Salmonella", Label: "Salmonella"},
		{XID: uuid.New(), CategoryID: 3, Value: "Appearance", Label: "Appearance"},
		{XID: uuid.New(), CategoryID: 3, Value: "Solubility", Label: "Solubility"},
		{XID: uuid.New(), CategoryID: 3, Value: "Particle Size", Label: "Particle Size"},
		{XID: uuid.New(), CategoryID: 2, Value: "Identification", Label: "Identification"},
		{XID: uuid.New(), CategoryID: 2, Value: "Assay", Label: "Assay"},
		{XID: uuid.New(), CategoryID: 2, Value: "Impurity", Label: "Impurity"},
		{XID: uuid.New(), CategoryID: 2, Value: "pH", Label: "pH"},
		{XID: uuid.New(), CategoryID: 2, Value: "Heavy Metals", Label: "Heavy Metals"},
		{XID: uuid.New(), CategoryID: 2, Value: "Residual Solvents", Label: "Residual Solvents"},
		{XID: uuid.New(), CategoryID: 2, Value: "Water Content", Label: "Water Content"},
	}
	for _, t := range tests {
		if err := q.SeedTestList(ctx, sqlc.SeedTestListParams{Xid: t.XID, CategoryID: t.CategoryID, Value: t.Value, Label: t.Label}); err != nil {
			log.Fatalf("Failed to seed test list data: %v", err)
		}
	}
}
