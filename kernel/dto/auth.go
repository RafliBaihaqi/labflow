package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name" binding:"required"`
	Company     string `json:"company"`
	Address     string `json:"address"`
}

type RegisterResponse struct {
	ID          int       `json:"id"`
	Xid         uuid.UUID `json:"xid"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Name        string    `json:"name"`
	Company     string    `json:"company,omitempty"`
	Address     string    `json:"address,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string    `json:"message"`
}

type LoginRequest struct {
	Identifier     string `json:"identifier" binding:"required"`
	Password       string `json:"password" binding:"required"`
	TurnstileToken string `json:"turnstileToken" binding:"required"`
}

type LoginResponse struct {
	Message string  `json:"message"`
	User    UserDTO `json:"user"`
}

type UserDTO struct {
	ID          int32     `json:"id"`
	Xid         uuid.UUID `json:"xid"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Company     string    `json:"company,omitempty"`
}

type VerifyEmail struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp"   binding:"required"`
}

type ResendOTP struct {
	Email string `json:"email" binding:"required,email"`
}

type MeResponse struct {
	ID          int32     `json:"id"`
	Xid         uuid.UUID `json:"xid"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Company     string    `json:"company,omitempty"`
	Address     string    `json:"address,omitempty"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	Avatar      *string   `json:"avatar"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyResetOtpRequest struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp" binding:"required"`
}

type VerifyResetOtpResponse struct {
	ResetToken string `json:"reset_token"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}
