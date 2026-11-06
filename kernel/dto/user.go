package dto

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	ID          int       `json:"id"`
	XID         uuid.UUID `json:"xid"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Name        string    `json:"name"`
	Company     string    `json:"company"`
	Address     string    `json:"address"`
	Avatar      *string   `json:"avatar"`
}

type UpdateUserRequest struct {
	PhoneNumber *string `json:"phone_number"`
	Company     *string `json:"company"`
	Address     *string `json:"address"`
	Avatar      *string `json:"avatar"`
}

type UpdateUserResponse struct {
	Xid         uuid.UUID `json:"xid"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Company     string    `json:"company,omitempty"`
	Address     string    `json:"address,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
	Message     string    `json:"message"`
}
