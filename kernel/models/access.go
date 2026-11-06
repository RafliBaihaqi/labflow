package models

type Access struct {
	ID         int     `gorm:"column:access_id;primaryKey" json:"access_id"`
	AccessName *string `gorm:"column:access_name;type:varchar(255)" json:"access_name"`
}
