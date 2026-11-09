package models

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID             int        `json:"id"`
	XID            uuid.UUID  `json:"xid"`
	ReservationID  int        `json:"reservation_id"`
	UserID         int        `json:"user_id"`
	CertificateURL string     `json:"certificate_url"`
	HasFeedback    *bool      `json:"has_feedback"`
	IssuedAt       *time.Time `json:"issued_at"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
