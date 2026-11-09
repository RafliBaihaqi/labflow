package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID                  int        `json:"id"`
	XID                 uuid.UUID  `json:"xid"`
	UserID              int        `json:"user_id"`
	ReservationDate     time.Time  `json:"reservation_date"`
	Status              *string    `json:"status"`
	SampleName          string     `json:"sample_name"`
	SampleAmount        float64    `json:"sample_amount"`
	ReceivedNote        *string    `json:"received_note,omitempty"`
	TestRequestNumber   *string    `json:"test_request_number,omitempty"`
	SampleResidue       string     `json:"sample_residue"`
	SpecialInstructions *string    `json:"special_instructions,omitempty"`
	POURL               *string    `json:"po_url,omitempty"`
	PONumber            *string    `json:"po_number,omitempty"`
	SampleType          string     `json:"sample_type"`
	FeedbackRequired    *bool      `json:"feedback_required"`
	MRApproved          *bool      `json:"mr_approved"`
	HeadApproved        *bool      `json:"head_approved"`
	MRApprovalDate      *time.Time `json:"mr_approval_date,omitempty"`
	HeadApprovalDate    *time.Time `json:"head_approval_date,omitempty"`
	SamplePhoto         *string    `json:"sample_photo,omitempty"`
	Note                *string    `json:"note,omitempty"`
	StartTestDate       *time.Time `json:"start_test_date,omitempty"`
	EndTestDate         *time.Time `json:"end_test_date,omitempty"`
	SAPNotes            *string    `json:"sap_notes,omitempty"`
	BatchNumber         *string    `json:"batch_number,omitempty"`
	LeadTime            *int       `json:"lead_time,omitempty"`
	SampleTakenDate     time.Time  `json:"sample_taken_date"`
	SampleTakenBy       string     `json:"sample_taken_by"`
	IsKan               *bool      `json:"is_kan"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	UpdateBy            *int       `json:"updated_by"`
	Tests               []Test     `json:"tests,omitempty"`
	User                User       `json:"user"`
}
