package repo

import (
	"context"
	"fmt"
	"labflow/models"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type TestListRepo struct {
	db *gorm.DB
}

func NewTestListRepo(db *gorm.DB) *TestListRepo {
	return &TestListRepo{db: db}
}

func (tl *TestListRepo) GetTestListByCategory(ctx context.Context, categoryID int) ([]models.TestList, error) {
	var testlist []models.TestList
	if err := tl.db.WithContext(ctx).Preload("Category").Where("category_id = ?", categoryID).Find(&testlist).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return testlist, nil
}

func (tl *TestListRepo) ListTestLists(ctx context.Context, search string, categoryID *int, page, limit int) ([]models.TestList, int64, error) {
	var testlists []models.TestList
	var total int64
	query := tl.db.WithContext(ctx).Model(&models.TestList{}).Preload("Category")
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("value ILIKE ? OR label ILIKE ?", like, like)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapDBError(err)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&testlists).Error; err != nil {
		return nil, 0, wrapDBError(err)
	}
	return testlists, total, nil
}

func (tl *TestListRepo) CreateTestList(ctx context.Context, input models.TestList) (*models.TestList, error) {
	if input.XID == uuid.Nil {
		input.XID = uuid.New()
	}
	if input.IsCurrent == nil {
		isCurrent := true
		input.IsCurrent = &isCurrent
	}
	if err := tl.db.WithContext(ctx).Create(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to create test list: %w", wrapDBError(err))
	}
	if err := tl.db.WithContext(ctx).Preload("Category").First(&input, input.ID).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &input, nil
}

func (tl *TestListRepo) GetTestListByXid(ctx context.Context, xid uuid.UUID) (*models.TestList, error) {
	var testList models.TestList
	if err := tl.db.WithContext(ctx).Preload("Category").Where("xid = ?", xid).First(&testList).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &testList, nil
}

func (tl *TestListRepo) CreateNewVersion(ctx context.Context, xid uuid.UUID, input models.TestList) (*models.TestList, error) {
	var created models.TestList
	err := tl.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.TestList{}).Where("xid = ?", xid).Update("is_current", false).Error; err != nil {
			return fmt.Errorf("failed to deactivate previous version: %w", wrapDBError(err))
		}
		input.ID = 0
		input.XID = uuid.New()
		if err := tx.Create(&input).Error; err != nil {
			return fmt.Errorf("failed to create new version: %w", wrapDBError(err))
		}
		if err := tx.Preload("Category").First(&created, input.ID).Error; err != nil {
			return wrapDBError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (tl *TestListRepo) DeactivateByXid(ctx context.Context, xid uuid.UUID) (*models.TestList, error) {
	if err := tl.db.WithContext(ctx).Model(&models.TestList{}).Where("xid = ?", xid).Update("is_current", false).Error; err != nil {
		return nil, fmt.Errorf("failed to deactivate test list: %w", wrapDBError(err))
	}
	return tl.GetTestListByXid(ctx, xid)
}
