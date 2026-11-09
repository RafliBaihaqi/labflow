package models

import (
	"time"

	"github.com/google/uuid"
)

type TestList struct {
	ID                 int        `json:"id"`
	XID                uuid.UUID  `json:"xid"`
	CategoryID         int        `json:"category_id"`
	Value              string     `json:"value"`
	Label              string     `json:"label"`
	SampleAmountLiquid *int       `json:"sample_amount_liquid,omitempty"`
	SampleAmountSolid  *int       `json:"sample_amount_solid,omitempty"`
	SampleAmountPowder *int       `json:"sample_amount_powder,omitempty"`
	Leadtime           *int       `json:"leadtime,omitempty"`
	Version            int        `json:"version"`
	IsCurrent          *bool      `json:"is_current"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	Category           Category   `json:"category"`
}
