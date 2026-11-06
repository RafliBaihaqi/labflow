package seeders

import (
	"labflow/models"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedRole(g *gorm.DB) {
	db := g
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
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error; err != nil {
		log.Fatalf("Failed to seed role data: %v", err)
	}
}
