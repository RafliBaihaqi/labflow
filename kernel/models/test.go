package models

import (
	"time"

	"github.com/google/uuid"
)

// TestItem is the lightweight parameter for creating test line items.
// It is used by both the service and repo layers.
type TestItem struct {
	TestListID int
	Method     string
	Criteria   string
	TestUnit   string
}

type Test struct {
	ID            int        `json:"id"`
	XID           uuid.UUID  `json:"xid"`
	ReservationID int        `json:"reservation_id"`
	TestListID    int        `json:"test_list_id"`
	Method        string     `json:"method"`
	Criteria      string     `json:"criteria"`
	Status        *string    `json:"status,omitempty"`
	TestUnit      string     `json:"test_unit"`
	TestResult    *string    `json:"test_result,omitempty"`
	TestSummary   *string    `json:"test_summary,omitempty"`
	ApprovedDate  *time.Time `json:"approved_date,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdateBy      *int       `json:"updated_by"`
	IsApproved    *bool      `json:"is_approved"`
	TestList      TestList   `json:"test_list"`
}
