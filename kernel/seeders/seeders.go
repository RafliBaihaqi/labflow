package seeders

import "gorm.io/gorm"

func RunSeeders(db *gorm.DB) {
	SeedAccess(db)
	SeedRole(db)
	SeedRoleAccess(db)
	SeedUser(db)
	SeedCategory(db)
	SeedTestList(db)
}
