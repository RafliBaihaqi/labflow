package service

import (
	"context"
	"labflow/models"

	"github.com/google/uuid"
)

type userRepo interface {
	UpdateUserByXid(ctx context.Context, xid uuid.UUID, phoneNumber, company, address, avatar *string, updatedBy int) (*models.User, error)
}

type UpdateUserParams struct {
	PhoneNumber *string
	Company     *string
	Address     *string
	Avatar      *string
}

type UserResult struct {
	Xid         uuid.UUID
	PhoneNumber string
	Company     string
	Address     string
	Avatar      *string
	Message     string
}

type UserService struct {
	repo userRepo
}

func NewUserService(repo userRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) UpdateUserByXid(ctx context.Context, updatedBy int, xid uuid.UUID, req UpdateUserParams) (UserResult, error) {
	user, err := us.repo.UpdateUserByXid(ctx, xid, req.PhoneNumber, req.Company, req.Address, req.Avatar, updatedBy)
	if err != nil {
		return UserResult{}, ErrFailedUpdate
	}

	return UserResult{
		Xid:         user.XID,
		PhoneNumber: user.PhoneNumber,
		Company:     user.Company,
		Address:     user.Address,
		Avatar:      FileURL(FileCategoryUser, user.Avatar),
		Message:     "User updated",
	}, nil
}
