package seeders

import (
	"context"
	"labflow/models"
	"labflow/sqlc"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SeedUser(ctx context.Context, pool *pgxpool.Pool) {
	q := sqlc.New(pool)
	defaultPassword := "password123"
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), 12)
	if err != nil {
		log.Fatalf("Failed to hash default password: %v", err)
	}
	hashedPassword := string(hashedPasswordBytes)
	users := []models.User{
		{
			XID:         uuid.New(),
			Username:    "superuser",
			Email:       "superuser@labflow.local",
			Name:        "Superuser",
			Password:    hashedPassword,
			PhoneNumber: "081111111111",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "admin",
			Email:       "admin@labflow.local",
			Name:        "Admin",
			Password:    hashedPassword,
			PhoneNumber: "081111111112",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "qa",
			Email:       "qa@labflow.local",
			Name:        "Quality Assurance",
			Password:    hashedPassword,
			PhoneNumber: "081111111113",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "customer",
			Email:       "customer@labflow.local",
			Name:        "Customer",
			Password:    hashedPassword,
			PhoneNumber: "081111111114",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "micro_analyst",
			Email:       "micro_analyst@labflow.local",
			Name:        "Micro Analyst",
			Password:    hashedPassword,
			PhoneNumber: "081111111115",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "chem_analyst",
			Email:       "chem_analyst@labflow.local",
			Name:        "Chemical Analyst",
			Password:    hashedPassword,
			PhoneNumber: "081111111116",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "micro_qc",
			Email:       "micro_qc@labflow.local",
			Name:        "Micro QC",
			Password:    hashedPassword,
			PhoneNumber: "081111111117",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "chem_qc",
			Email:       "chem_qc@labflow.local",
			Name:        "Chemical QC",
			Password:    hashedPassword,
			PhoneNumber: "081111111118",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "lab_head",
			Email:       "lab_head@labflow.local",
			Name:        "Lab Head",
			Password:    hashedPassword,
			PhoneNumber: "081111111119",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
		{
			XID:         uuid.New(),
			Username:    "mr",
			Email:       "mr@labflow.local",
			Name:        "Management Representative",
			Password:    hashedPassword,
			PhoneNumber: "081111111120",
			Address:     "Alamat 1",
			Company:     "Company 1",
			IsVerified:  func(b bool) *bool { return &b }(true),
		},
	}

	for _, u := range users {
		// Upsert user
		if err := q.SeedUser(ctx, sqlc.SeedUserParams{
			Xid:         u.XID,
			Username:    u.Username,
			Password:    u.Password,
			Email:       u.Email,
			PhoneNumber: &u.PhoneNumber,
			Name:        u.Name,
			Company:     &u.Company,
			Address:     &u.Address,
			IsVerified:  u.IsVerified,
		}); err != nil {
			log.Fatalf("Failed to seed user %s: %v", u.Username, err)
		}
		// Fetch back to get generated ID
		existing, err := q.GetUserByUsername(ctx, u.Username)
		if err != nil {
			log.Printf("User %s already exists or not found, skipping role assignment", u.Username)
			continue
		}
		// Map roles
		var roleID int
		switch existing.Username {
		case "superuser", "admin":
			roleID = 1
		case "customer":
			roleID = 2
		case "qa":
			roleID = 3
		case "micro_analyst":
			roleID = 4
		case "chem_analyst":
			roleID = 5
		case "micro_qc":
			roleID = 6
		case "chem_qc":
			roleID = 7
		case "lab_head":
			roleID = 8
		case "mr":
			roleID = 9
		}
		// Insert user_role join
		if roleID > 0 {
			if err := q.SeedUserRole(ctx, sqlc.SeedUserRoleParams{UserID: existing.ID, RoleID: roleID}); err != nil {
				log.Fatalf("Failed to assign role %d to user %s: %v", roleID, existing.Username, err)
			}
		}
		log.Printf("Seeded user: %s (%s) with role_id %d", existing.Username, existing.Email, roleID)
	}
}
