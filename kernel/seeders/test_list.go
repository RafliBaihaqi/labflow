package seeders

import (
	"labflow/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedTestList(g *gorm.DB) {
	db := g
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
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&tests).Error; err != nil {
		log.Fatalf("Failed to seed test list data: %v", err)
	}
}
