package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestRepo struct {
	db *gorm.DB
}

func NewTestRepo(db *gorm.DB) *TestRepo {
	return &TestRepo{db: db}
}

type UpdateTestParams struct {
	TestResult  *string
	TestSummary *string
	Notes       *string
	UpdatedBy   *int
	IsApproved  *bool
}

func (t *TestRepo) GetTestByXid(ctx context.Context, xid uuid.UUID) (*models.Test, error) {
	var test *models.Test
	if err := t.db.WithContext(ctx).Preload("TestList.Category").Where("xid = ?", xid).First(&test).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return test, nil
}

func (t *TestRepo) UpdateTestByXid(ctx context.Context, xid uuid.UUID, testResult, testSummary, notes *string, updatedBy int) (*models.Test, error) {
	params := UpdateTestParams{
		TestResult:  testResult,
		TestSummary: testSummary,
		Notes:       notes,
		UpdatedBy:   &updatedBy,
	}

	test := map[string]interface{}{
		"test_result":  params.TestResult,
		"test_summary": params.TestSummary,
		"notes":        params.Notes,
		"updated_by":   params.UpdatedBy,
	}

	if err := t.db.Model(&models.Test{}).Where("xid = ?", xid).Updates(test).Error; err != nil {
		return nil, fmt.Errorf("failed to update test: %w", wrapDBError(err))
	}

	var updateTest models.Test
	if err := t.db.Where("xid = ?", xid).First(&updateTest).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &updateTest, nil
}

func (t *TestRepo) UpdateIsApproved(ctx context.Context, xid uuid.UUID, isApproved *bool, updatedBy int) (*models.Test, error) {
	params := UpdateTestParams{
		IsApproved: isApproved,
		UpdatedBy:  &updatedBy,
	}

	test := map[string]interface{}{
		"approved_date": time.Now(),
		"updated_by":    params.UpdatedBy,
		"is_approved":   params.IsApproved,
	}

	if err := t.db.Model(&models.Test{}).Where("xid = ?", xid).Updates(test).Error; err != nil {
		return nil, fmt.Errorf("failed to update test: %w", wrapDBError(err))
	}

	var ad models.Test
	if err := t.db.Where("xid = ?", xid).First(&ad).Error; err != nil {
		return nil, wrapDBError(err)
	}

	return &ad, nil
}
