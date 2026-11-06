package models

import (
	"time"

	"github.com/google/uuid"
)

type TestList struct {
	ID                 int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	XID                uuid.UUID  `gorm:"column:xid;type:uuid;not null" json:"xid"`
	CategoryID         int        `gorm:"column:category_id;type:int;not null;index" json:"category_id"`
	Value              string     `gorm:"column:value;type:varchar(255);not null" json:"value"`
	Label              string     `gorm:"column:label;type:varchar(255);not null" json:"label"`
	SampleAmountLiquid *int       `gorm:"column:sample_amount_liquid;type:int" json:"sample_amount_liquid,omitempty"`
	SampleAmountSolid  *int       `gorm:"column:sample_amount_solid;type:int" json:"sample_amount_solid,omitempty"`
	SampleAmountPowder *int       `gorm:"column:sample_amount_powder;type:int" json:"sample_amount_powder,omitempty"`
	Leadtime           *int       `gorm:"column:leadtime;type:int" json:"leadtime,omitempty"`
	Version            int        `gorm:"column:version;type:int;not null;default:1" json:"version"`
	IsCurrent          *bool      `gorm:"column:is_current;type:boolean;default:true" json:"is_current"`
	CreatedAt          *time.Time `gorm:"column:created_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          *time.Time `gorm:"column:updated_at;type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Category           Category   `gorm:"foreignKey:CategoryID" json:"category"`
}
