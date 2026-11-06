package models

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID             int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	XID            uuid.UUID  `gorm:"column:xid;type:uuid;not null;uniqueIndex" json:"xid"`
	ReservationID  int        `gorm:"column:reservation_id;type:int;not null;index" json:"reservation_id"`
	UserID         int        `gorm:"column:user_id;type:int;not null;index" json:"user_id"`
	CertificateURL string     `gorm:"column:certificate_url;type:text;not null" json:"certificate_url"`
	HasFeedback    *bool      `gorm:"column:has_feedback;type:boolean;default:false" json:"has_feedback"`
	IssuedAt       *time.Time `gorm:"column:issued_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"issued_at"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
