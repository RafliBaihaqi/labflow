package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTestRequest struct {
	TestListID int    `json:"test_list_id" binding:"required"`
	Method     string `json:"method" binding:"required"`
	Criteria   string `json:"criteria" binding:"required"`
	TestUnit   string `json:"test_unit" binding:"required"`
}

type UpdateTestRequest struct {
	TestResult  *string `json:"test_result"`
	TestSummary *string `json:"test_summary"`
	Notes       *string `json:"notes"`
	IsApproved  *bool   `json:"is_approved"`
}
type TestResponse struct {
	ID            int        `json:"id"`
	XID           uuid.UUID  `json:"xid"`
	ReservationID int        `json:"reservation_id"`
	TestListId    int        `json:"test_list_id"`
	Method        string     `json:"method"`
	Criteria      string     `json:"criteria"`
	Status        *string    `json:"status"`
	TestUnit      string     `json:"test_unit"`
	TestResult    *string    `json:"test_result"`
	TestSummary   *string    `json:"test_summary"`
	ApprovedDate  *time.Time `json:"approved_date"`
	Notes         *string    `json:"notes"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	TestList      string     `json:"test_list_name"`
	CategoryName  string     `json:"category_name"`
}

type TestCatalogResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type TestListInput struct {
	CategoryID         int    `json:"category_id" binding:"required"`
	Value              string `json:"value" binding:"required"`
	Label              string `json:"label" binding:"required"`
	SampleAmountLiquid *int   `json:"sample_amount_liquid"`
	SampleAmountSolid  *int   `json:"sample_amount_solid"`
	SampleAmountPowder *int   `json:"sample_amount_powder"`
	Leadtime           *int   `json:"leadtime"`
}

type TestListResponse struct {
	XID                string `json:"xid"`
	CategoryID         int    `json:"category_id"`
	Value              string `json:"value"`
	Label              string `json:"label"`
	SampleAmountLiquid *int   `json:"sample_amount_liquid,omitempty"`
	SampleAmountSolid  *int   `json:"sample_amount_solid,omitempty"`
	SampleAmountPowder *int   `json:"sample_amount_powder,omitempty"`
	Leadtime           *int   `json:"leadtime,omitempty"`
	Version            int    `json:"version"`
	IsCurrent          *bool  `json:"is_current,omitempty"`
	CategoryName       string `json:"category_name"`
}

type PaginatedTestListResponse struct {
	Data       []TestListResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}
