package models

import (
	"time"

	"github.com/google/uuid"
)

type Quotation struct {
	ID          int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	XID         uuid.UUID  `gorm:"column:xid;type:uuid;not null;uniqueIndex" json:"xid"`
	Email       string     `gorm:"column:email;type:text;not null" json:"email"`
	CompanyName string     `gorm:"column:company_name;type:text;not null" json:"company_name"`
	PhoneNumber *string    `gorm:"column:phone_number;type:text" json:"phone_number,omitempty"`
	Message     *string    `gorm:"column:message;type:text" json:"message,omitempty"`
	Status      *string    `gorm:"column:status;type:text;default:'pending'" json:"status"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
