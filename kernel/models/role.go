package models

type Role struct {
	ID       int      `gorm:"column:role_id;primaryKey" json:"role_id"`
	RoleName *string  `gorm:"column:role_name;type:varchar(255)" json:"role_name"`
	Accesses []Access `gorm:"many2many:role_accesses;joinForeignKey:role_id;joinReferences:access_id" json:"accesses,omitempty"`
}
