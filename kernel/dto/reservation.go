package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateReservationRequest struct {
	SampleName          string              `json:"sample_name" binding:"required"`
	SampleAmount        float64             `json:"sample_amount" binding:"required"`
	ReceivedNote        *string             `json:"received_note"`
	SampleResidue       string              `json:"sample_residue"`
	SpecialInstructions *string             `json:"special_instructions" binding:"required"`
	SampleType          string              `json:"sample_type" binding:"required"`
	Note                *string             `json:"note"`
	BatchNumber         *string             `json:"batch_number"`
	SampleTakenDate     time.Time           `json:"sample_taken_date" binding:"required"`
	SampleTakenBy       string              `json:"sample_taken_by" binding:"required"`
	Tests               []CreateTestRequest `json:"tests" binding:"required"`
}

type UpdateReservationRequest struct {
	POUrl         *string    `json:"po_url"`
	PONumber      *string    `json:"po_number"`
	StartTestDate *time.Time `json:"start_test_date"`
	LeadTime      *int       `json:"lead_time"`
	SamplePhoto   *string    `json:"sample_photo"`
}

type UpdateStatusReq struct {
	Status      string  `json:"status" binding:"required"`
	SamplePhoto *string `json:"sample_photo,omitempty"`
}

type CreateReservationResponse struct {
	ID                  int                 `json:"id"`
	Xid                 uuid.UUID           `json:"xid"`
	SampleName          string              `json:"sample_name"`
	SampleAmount        float64             `json:"sample_amount"`
	SpecialInstructions *string             `json:"special_instructions"`
	SampleType          string              `json:"sample_type"`
	BatchNumber         *string             `json:"batch_number"`
	SampleTakenDate     time.Time           `json:"sample_taken_date"`
	SampleTakenBy       string              `json:"sample_taken_by"`
	TestRequestNumber   *string             `json:"test_request_number"`
	Message             string              `json:"message"`
	Tests               []CreateTestRequest `json:"tests"`
}

type UploadResponse struct {
	Message    string `json:"message"`
	ObjectName string `json:"object_name"`
}

type ReservationStatsResponse struct {
	Total          int64 `json:"total"`
	Pending        int64 `json:"pending"`
	SampleReceived int64 `json:"sample_received"`
	TestOnProgress int64 `json:"test_on_progress"`
	TestCompleted  int64 `json:"test_completed"`
	Rejected       int64 `json:"rejected"`
}

type ReservationResponse struct {
	ID                  int            `json:"id"`
	Xid                 uuid.UUID      `json:"xid"`
	UserID              int            `json:"user_id"`
	ReservationDate     time.Time      `json:"reservation_date"`
	Status              *string        `json:"status"`
	SampleName          string         `json:"sample_name"`
	SampleAmount        float64        `json:"sample_amount"`
	ReceivedNote        *string        `json:"received_note"`
	TestRequestNumber   *string        `json:"test_request_number"`
	SampleResidue       string         `json:"sample_residue"`
	SpecialInstructions *string        `json:"special_instructions"`
	POURL               *string        `json:"po_url"`
	PONumber            *string        `json:"po_number"`
	SampleType          string         `json:"sample_type"`
	FeedbackRequired    *bool          `json:"feedback_required"`
	MRApproved          *bool          `json:"mr_approved"`
	HeadApproved        *bool          `json:"head_approved"`
	MRApprovalDate      *time.Time     `json:"mr_approval_date"`
	HeadApprovalDate    *time.Time     `json:"head_approval_date"`
	SamplePhoto         *string        `json:"sample_photo"`
	Note                *string        `json:"note"`
	StartTestDate       *time.Time     `json:"start_test_date"`
	EndTestDate         *time.Time     `json:"end_test_date"`
	SAPNotes            *string        `json:"sap_notes"`
	BatchNumber         *string        `json:"batch_number"`
	LeadTime            *int           `json:"lead_time"`
	SampleTakenDate     time.Time      `json:"sample_taken_date"`
	SampleTakenBy       string         `json:"sample_taken_by"`
	IsKan               *bool          `json:"is_kan"`
	CreatedAt           *time.Time     `json:"created_at"`
	UpdatedAt           *time.Time     `json:"updated_at"`
	UpdatedBy           *int           `json:"updated_by"`
	Message             string         `json:"message"`
	Tests               []TestResponse `json:"tests"`
	User                UserResponse   `json:"user"`
}
