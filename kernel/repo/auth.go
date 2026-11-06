package repo

import (
	"context"
	"fmt"
	"labflow/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

type CreateUserParams struct {
	Username    string
	Password    string
	Email       string
	PhoneNumber string
	Name        string
	Company     string
	Address     string
	RoleID      int
}

func (r *AuthRepo) CreateUser(ctx context.Context, username, password, email string, phoneNumber string, name string, company, address string, roleID int) (*models.User, error) {
	params := CreateUserParams{
		Username:    username,
		Password:    password,
		Email:       email,
		PhoneNumber: phoneNumber,
		Name:        name,
		Company:     company,
		Address:     address,
		RoleID:      roleID,
	}

	user := models.User{
		XID:         uuid.New(),
		Username:    params.Username,
		Password:    params.Password,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		Name:        params.Name,
		Company:     params.Company,
		Address:     params.Address,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to insert user: %w", wrapDBError(err))
		}

		if params.RoleID > 0 {
			userRole := models.UserRole{
				UserID: user.ID,
				RoleID: params.RoleID,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				return fmt.Errorf("failed to map user role: %w", wrapDBError(err))
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username=?", username).First(&user).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &user, nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email=?", email).First(&user).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &user, nil
}

func (r *AuthRepo) UpdateUserVerifiedStatus(ctx context.Context, email string) (*models.User, error) {
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Update("is_verified", true).Error; err != nil {
		return nil, wrapDBError(err)
	}
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &user, nil
}

func (r *AuthRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles.Accesses").First(&user, userID).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &user, nil
}

func (r *AuthRepo) UpdateUserPassword(ctx context.Context, email string, hashedPassword string) (*models.User, error) {
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Update("password", hashedPassword).Error; err != nil {
		return nil, wrapDBError(err)
	}
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &user, nil
}
