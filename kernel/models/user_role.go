package models

type UserRole struct {
	UserID int `gorm:"column:user_id;primaryKey" json:"user_id"`
	RoleID int `gorm:"column:role_id;primaryKey" json:"role_id"`
}
