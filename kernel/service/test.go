package service

import (
	"context"
	"errors"
	"labflow/models"

	"github.com/google/uuid"
)

var ErrUnauthorizedRole = errors.New("unauthorized role for this test category")

type testRepo interface {
	GetTestByXid(ctx context.Context, xid uuid.UUID) (*models.Test, error)
	UpdateTestByXid(ctx context.Context, xid uuid.UUID, testResult, testSummary, notes *string, updatedBy int) (*models.Test, error)
	UpdateIsApproved(ctx context.Context, xid uuid.UUID, isApproved *bool, updatedBy int) (*models.Test, error)
}

type roleReader interface {
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
}

type TestService struct {
	repo  testRepo
	roles roleReader
}

func NewTestService(repo testRepo, roles roleReader) *TestService {
	return &TestService{
		repo:  repo,
		roles: roles,
	}
}

type UpdateTestParams struct {
	TestResult  *string
	TestSummary *string
	Notes       *string
	UpdatedBy   *int
	IsApproved  *bool
}

func (ts *TestService) UpdateTestByXid(ctx context.Context, xid uuid.UUID, updatedBy int, req UpdateTestParams) (MessageResult, error) {
	test, err := ts.repo.GetTestByXid(ctx, xid)
	if err != nil {
		return MessageResult{}, ErrFailedFetch
	}

	user, err := ts.roles.GetUserByID(ctx, updatedBy)
	if err != nil {
		return MessageResult{}, err
	}

	categoryName := test.TestList.Category.Name
	if !userHasTestCategoryAccess(user.Roles, categoryName) {
		return MessageResult{}, ErrUnauthorizedRole
	}

	_, err = ts.repo.UpdateTestByXid(ctx, xid, req.TestResult, req.TestSummary, req.Notes, updatedBy)
	if err != nil {
		return MessageResult{}, ErrFailedUpdate
	}

	return MessageResult{
		Message: "Test updated",
	}, nil
}

func (ts *TestService) UpdateIsApproved(ctx context.Context, xid uuid.UUID, updatedBy int, req UpdateTestParams) (MessageResult, error) {
	_, err := ts.repo.UpdateIsApproved(ctx, xid, req.IsApproved, updatedBy)
	if err != nil {
		return MessageResult{}, ErrFailedUpdate
	}

	return MessageResult{
		Message: "Test updated",
	}, nil
}

func userHasTestCategoryAccess(roles []models.Role, categoryName string) bool {
	for _, role := range roles {
		if role.RoleName == nil {
			continue
		}
		switch *role.RoleName {
		case "admin":
			return true
		case "micro-analyst", "micro-qc":
			if categoryName == "Microbiology" {
				return true
			}
		case "chemical-analyst", "chemical-qc":
			if categoryName == "Chemistry" {
				return true
			}
		}
	}
	return false
}
