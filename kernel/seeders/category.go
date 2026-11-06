package seeders

import (
	"labflow/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedCategory(g *gorm.DB) {
	db := g
	categories := []models.Category{
		{XID: uuid.New(), Name: "Microbiology", Label: func(s string) *string { return &s }("Microbiology")},
		{XID: uuid.New(), Name: "Chemistry", Label: func(s string) *string { return &s }("Chemistry")},
		{XID: uuid.New(), Name: "Physical", Label: func(s string) *string { return &s }("Physical")},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&categories).Error; err != nil {
		log.Fatalf("Failed to seed category data: %v", err)
	}
}
