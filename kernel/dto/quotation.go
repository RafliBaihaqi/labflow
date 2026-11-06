package dto

import "github.com/google/uuid"

type CreateQuotationRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	CompanyName string  `json:"company_name" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Message     *string `json:"message,omitempty"`
}

type CreateQuotationResponse struct {
	XID         uuid.UUID `json:"xid"`
	Email       string    `json:"email"`
	CompanyName string    `json:"company_name"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Message     *string   `json:"message,omitempty"`
}

type UpdateQuotationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type QuotationResponse struct {
	XID         uuid.UUID `json:"xid"`
	Email       string    `json:"email"`
	CompanyName string    `json:"company_name"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Message     *string   `json:"message,omitempty"`
	Status      *string   `json:"status,omitempty"`
}

type PaginatedQuotationResponse struct {
	Data       []QuotationResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}
