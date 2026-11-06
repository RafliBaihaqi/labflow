package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

type UpdateUserParams struct {
	PhoneNumber *string
	Company     *string
	Address     *string
	Avatar      *string
	UpdatedBy   *int
}

func (u *UserRepo) UpdateUserByXid(ctx context.Context, xid uuid.UUID, phoneNumber, company, address, avatar *string, updatedBy int) (*models.User, error) {
	params := UpdateUserParams{
		PhoneNumber: phoneNumber,
		Company:     company,
		Address:     address,
		Avatar:      avatar,
		UpdatedBy:   &updatedBy,
	}

	updates := map[string]interface{}{
		"phone_number": params.PhoneNumber,
		"company":      params.Company,
		"address":      params.Address,
		"avatar":       params.Avatar,
		"updated_at":   time.Now(),
		"updated_by":   params.UpdatedBy,
	}

	if err := u.db.Model(&models.User{}).Where("xid = ?", xid).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update reservation: %w", wrapDBError(err))
	}

	var updated models.User
	if err := u.db.Where("xid = ?", xid).First(&updated).Error; err != nil {
		return nil, wrapDBError(err)
	}

	return &updated, nil
}
