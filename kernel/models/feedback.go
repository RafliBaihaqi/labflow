package models

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID               int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	XID              uuid.UUID  `gorm:"column:xid;type:uuid;not null;uniqueIndex" json:"xid"`
	UserID           int        `gorm:"column:user_id;type:int;not null;index" json:"user_id"`
	CertificateID    int        `gorm:"column:certificate_id;type:int;not null;index" json:"certificate_id"`
	ResultClarity    *int       `gorm:"column:result_clarity;type:int" json:"result_clarity,omitempty"`
	Communication    *int       `gorm:"column:communication;type:int" json:"communication,omitempty"`
	CustomerService  *int       `gorm:"column:customer_service;type:int" json:"customer_service,omitempty"`
	RecommendService *int       `gorm:"column:reccomend_service;type:int" json:"reccomend_service,omitempty"`
	ResultAccuracy   *int       `gorm:"column:result_accuracy;type:int" json:"result_accuracy,omitempty"`
	Notes            *string    `gorm:"column:notes;type:text" json:"notes,omitempty"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
