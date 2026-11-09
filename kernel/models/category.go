package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        int        `json:"id"`
	XID       uuid.UUID  `json:"xid"`
	Name      string     `json:"name"`
	Label     *string    `json:"label,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
