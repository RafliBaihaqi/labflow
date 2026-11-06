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
	ID            int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	XID           uuid.UUID  `gorm:"column:xid;type:uuid;not null;uniqueIndex" json:"xid"`
	ReservationID int        `gorm:"column:reservation_id;type:int;not null;index" json:"reservation_id"`
	TestListID    int        `gorm:"column:test_list_id;type:int;not null;index" json:"test_list_id"`
	Method        string     `gorm:"column:method;type:varchar(255);not null" json:"method"`
	Criteria      string     `gorm:"column:criteria;type:varchar(255);not null" json:"criteria"`
	Status        *string    `gorm:"column:status;type:varchar(255)" json:"status,omitempty"`
	TestUnit      string     `gorm:"column:test_unit;type:varchar(255);not null" json:"test_unit"`
	TestResult    *string    `gorm:"column:test_result;type:varchar(255)" json:"test_result,omitempty"`
	TestSummary   *string    `gorm:"column:test_summary;type:varchar(255)" json:"test_summary,omitempty"`
	ApprovedDate  *time.Time `gorm:"column:approved_date;type:timestamptz" json:"approved_date,omitempty"`
	Notes         *string    `gorm:"column:notes;type:varchar(500)" json:"notes,omitempty"`
	CreatedAt     *time.Time `gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdateBy      *int       `gorm:"column:updated_by;type:int" json:"updated_by"`
	IsApproved    *bool      `gorm:"column:is_approved;type:boolean;default:false" json:"is_approved"`
	TestList      TestList   `gorm:"foreignKey:TestListID" json:"test_list"`
}
