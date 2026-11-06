package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	XID       uuid.UUID  `gorm:"column:xid;type:uuid;not null;uniqueIndex" json:"xid"`
	Name      string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Label     *string    `gorm:"column:label;type:text" json:"label,omitempty"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
