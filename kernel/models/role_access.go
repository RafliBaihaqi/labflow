package models

type RoleAccess struct {
	RoleID   int `gorm:"column:role_id;primaryKey" json:"role_id"`
	AccessID int `gorm:"column:access_id;primaryKey" json:"access_id"`
}
