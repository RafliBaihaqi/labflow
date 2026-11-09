package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 int        `json:"id"`
	XID                uuid.UUID  `json:"xid"`
	Username           string     `json:"username"`
	Password           string     `json:"-"`
	Email              string     `json:"email"`
	PhoneNumber        string     `json:"phone_number,omitempty"`
	Name               string     `json:"name"`
	Company            string     `json:"company,omitempty"`
	Address            string     `json:"address,omitempty"`
	VerificationToken  *string    `json:"-"`
	VerificationExp    *time.Time `json:"-"`
	IsVerified         *bool      `json:"is_verified"`
	ResetPasswordToken *string    `json:"-"`
	ResetPasswordExp   *time.Time `json:"-"`
	SessionID          *string    `json:"-"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	UpdatedBy          *int       `json:"updated_by,omitempty"`
	Avatar             *string    `json:"avatar,omitempty"`
	Roles              []Role     `json:"roles,omitempty"`
}
