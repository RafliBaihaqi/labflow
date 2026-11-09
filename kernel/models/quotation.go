package models

import (
	"time"

	"github.com/google/uuid"
)

type Quotation struct {
	ID          int        `json:"id"`
	XID         uuid.UUID  `json:"xid"`
	Email       string     `json:"email"`
	CompanyName string     `json:"company_name"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	Message     *string    `json:"message,omitempty"`
	Status      *string    `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
