package models

type Role struct {
	ID       int      `json:"role_id"`
	RoleName *string  `json:"role_name"`
	Accesses []Access `json:"accesses,omitempty"`
}
