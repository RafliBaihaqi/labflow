package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID                  int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	XID                 uuid.UUID  `gorm:"column:xid;type:uuid;not null;uniqueIndex" json:"xid"`
	UserID              int        `gorm:"column:user_id;type:int;not null;index" json:"user_id"`
	ReservationDate     time.Time  `gorm:"column:reservation_date;type:date;not null" json:"reservation_date"`
	Status              *string    `gorm:"column:status;type:text;default:'pending'" json:"status"`
	SampleName          string     `gorm:"column:sample_name;type:text;not null" json:"sample_name"`
	SampleAmount        float64    `gorm:"column:sample_amount;type:double precision;not null" json:"sample_amount"`
	ReceivedNote        *string    `gorm:"column:received_note;type:text" json:"received_note,omitempty"`
	TestRequestNumber   *string    `gorm:"column:test_request_number;type:text" json:"test_request_number,omitempty"`
	SampleResidue       string     `gorm:"column:sample_residue;type:text;not null" json:"sample_residue"`
	SpecialInstructions *string    `gorm:"column:special_instructions;type:text" json:"special_instructions,omitempty"`
	POURL               *string    `gorm:"column:po_url;type:text" json:"po_url,omitempty"`
	PONumber            *string    `gorm:"column:po_number;type:text" json:"po_number,omitempty"`
	SampleType          string     `gorm:"column:sample_type;type:text;not null" json:"sample_type"`
	FeedbackRequired    *bool      `gorm:"column:feedback_required;type:boolean;default:true" json:"feedback_required"`
	MRApproved          *bool      `gorm:"column:mr_approved;type:boolean;default:false" json:"mr_approved"`
	HeadApproved        *bool      `gorm:"column:head_approved;type:boolean;default:false" json:"head_approved"`
	MRApprovalDate      *time.Time `gorm:"column:mr_approval_date;type:timestamptz" json:"mr_approval_date,omitempty"`
	HeadApprovalDate    *time.Time `gorm:"column:head_approval_date;type:timestamptz" json:"head_approval_date,omitempty"`
	SamplePhoto         *string    `gorm:"column:sample_photo;type:text" json:"sample_photo,omitempty"`
	Note                *string    `gorm:"column:note;type:text" json:"note,omitempty"`
	StartTestDate       *time.Time `gorm:"column:start_test_date;type:timestamptz" json:"start_test_date,omitempty"`
	EndTestDate         *time.Time `gorm:"column:end_test_date;type:timestamptz" json:"end_test_date,omitempty"`
	SAPNotes            *string    `gorm:"column:sap_notes;type:text" json:"sap_notes,omitempty"`
	BatchNumber         *string    `gorm:"column:batch_number;type:text;default:'-'" json:"batch_number,omitempty"`
	LeadTime            *int       `gorm:"column:lead_time;type:int" json:"lead_time,omitempty"`
	SampleTakenDate     time.Time  `gorm:"column:sample_taken_date;type:timestamptz;not null" json:"sample_taken_date"`
	SampleTakenBy       string     `gorm:"column:sample_taken_by;type:text;not null" json:"sample_taken_by"`
	IsKan               *bool      `gorm:"column:is_kan;type:boolean;default:false" json:"is_kan"`
	CreatedAt           *time.Time `gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdateBy            *int       `gorm:"column:updated_by;type:int" json:"updated_by"`
	Tests               []Test     `gorm:"foreignKey:ReservationID" json:"tests,omitempty"`
	User                User       `gorm:"foreignKey:UserID" json:"user"`
}
