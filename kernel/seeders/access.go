package seeders

import (
	"labflow/models"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAccess(g *gorm.DB) {
	db := g
	accesses := []models.Access{
		{ID: 1, AccessName: func(s string) *string { return &s }("create_user")},
		{ID: 2, AccessName: func(s string) *string { return &s }("read_user")},
		{ID: 3, AccessName: func(s string) *string { return &s }("update_user")},
		{ID: 4, AccessName: func(s string) *string { return &s }("delete_user")},
		{ID: 5, AccessName: func(s string) *string { return &s }("create_reservation")},
		{ID: 6, AccessName: func(s string) *string { return &s }("read_reservation")},
		{ID: 7, AccessName: func(s string) *string { return &s }("update_reservation")},
		{ID: 8, AccessName: func(s string) *string { return &s }("delete_reservation")},
		{ID: 9, AccessName: func(s string) *string { return &s }("create_test")},
		{ID: 10, AccessName: func(s string) *string { return &s }("read_test")},
		{ID: 11, AccessName: func(s string) *string { return &s }("update_test")},
		{ID: 12, AccessName: func(s string) *string { return &s }("delete_test")},
		{ID: 13, AccessName: func(s string) *string { return &s }("create_certificate")},
		{ID: 14, AccessName: func(s string) *string { return &s }("read_certificate")},
		{ID: 15, AccessName: func(s string) *string { return &s }("update_certificate")},
		{ID: 16, AccessName: func(s string) *string { return &s }("delete_certificate")},
		{ID: 17, AccessName: func(s string) *string { return &s }("update_reservation_status")},
		{ID: 18, AccessName: func(s string) *string { return &s }("update_user_avatar")},
		{ID: 19, AccessName: func(s string) *string { return &s }("read_quotation")},
		{ID: 20, AccessName: func(s string) *string { return &s }("update_quotation")},
		{ID: 21, AccessName: func(s string) *string { return &s }("read_dashboard")},
		{ID: 22, AccessName: func(s string) *string { return &s }("read_chat")},
		{ID: 23, AccessName: func(s string) *string { return &s }("create_chat")},
		{ID: 24, AccessName: func(s string) *string { return &s }("create_feedback")},
		{ID: 25, AccessName: func(s string) *string { return &s }("read_feedback")},
		{ID: 26, AccessName: func(s string) *string { return &s }("read_customer")},
		{ID: 27, AccessName: func(s string) *string { return &s }("export_reservation")},
		{ID: 28, AccessName: func(s string) *string { return &s }("read_lifecycle")},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&accesses).Error; err != nil {
		log.Fatalf("Failed to seed access data: %v", err)
	}
}
