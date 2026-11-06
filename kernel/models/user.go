package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	XID                uuid.UUID  `gorm:"column:xid;type:uuid;not null;uniqueIndex" json:"xid"`
	Username           string     `gorm:"column:username;type:varchar(255);not null;uniqueIndex" json:"username"`
	Password           string     `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Email              string     `gorm:"column:email;type:varchar(255);not null;uniqueIndex" json:"email"`
	PhoneNumber        string     `gorm:"column:phone_number;type:varchar(255);not null" json:"phone_number,omitempty"`
	Name               string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Company            string     `gorm:"column:company;type:varchar(255);not null" json:"company,omitempty"`
	Address            string     `gorm:"column:address;type:text;not null" json:"address,omitempty"`
	VerificationToken  *string    `gorm:"column:verification_token;type:text" json:"-"`
	VerificationExp    *time.Time `gorm:"column:verification_exp;type:timestamptz" json:"-"`
	IsVerified         *bool      `gorm:"column:is_verified;type:boolean;default:false" json:"is_verified"`
	ResetPasswordToken *string    `gorm:"column:reset_password_token;type:text" json:"-"`
	ResetPasswordExp   *time.Time `gorm:"column:reset_password_exp;type:timestamptz" json:"-"`
	SessionID          *string    `gorm:"column:session_id;type:varchar(255)" json:"-"`
	CreatedAt          *time.Time `gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          *time.Time `gorm:"column:updated_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy          *int       `gorm:"column:updated_by;type:int" json:"updated_by,omitempty"`
	Avatar             *string    `gorm:"column:avatar;type:string" json:"avatar,omitempty"`
	Roles              []Role     `gorm:"many2many:user_roles;joinForeignKey:user_id;joinReferences:role_id" json:"roles,omitempty"`
}
