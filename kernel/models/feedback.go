package models

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID               int        `json:"id"`
	XID              uuid.UUID  `json:"xid"`
	UserID           int        `json:"user_id"`
	CertificateID    int        `json:"certificate_id"`
	ResultClarity    *int       `json:"result_clarity,omitempty"`
	Communication    *int       `json:"communication,omitempty"`
	CustomerService  *int       `json:"customer_service,omitempty"`
	RecommendService *int       `json:"reccomend_service,omitempty"`
	ResultAccuracy   *int       `json:"result_accuracy,omitempty"`
	Notes            *string    `json:"notes,omitempty"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}
