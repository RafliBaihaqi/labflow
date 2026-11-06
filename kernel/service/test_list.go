package service

import (
	"context"
	"errors"
	"labflow/models"
	"strings"

	"github.com/google/uuid"
)

type testListRepo interface {
	GetTestListByCategory(ctx context.Context, categoryID int) ([]models.TestList, error)
	ListTestLists(ctx context.Context, search string, categoryID *int, page, limit int) ([]models.TestList, int64, error)
	CreateTestList(ctx context.Context, input models.TestList) (*models.TestList, error)
	GetTestListByXid(ctx context.Context, xid uuid.UUID) (*models.TestList, error)
	CreateNewVersion(ctx context.Context, xid uuid.UUID, input models.TestList) (*models.TestList, error)
	DeactivateByXid(ctx context.Context, xid uuid.UUID) (*models.TestList, error)
}

type TestListService struct {
	repo testListRepo
}

func NewTestListService(repo testListRepo) *TestListService {
	return &TestListService{
		repo: repo,
	}
}

type TestCatalogResult struct {
	ID       int
	Name     string
	Category string
}

var ErrInvalidTestList = errors.New("invalid test list")

type TestListResult struct {
	Data       []models.TestList
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

func (ts *TestListService) GetTestListByCategory(ctx context.Context, categoryID int) ([]TestCatalogResult, error) {
	tl, err := ts.repo.GetTestListByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	result := make([]TestCatalogResult, 0, len(tl))
	for _, t := range tl {
		result = append(result, TestCatalogResult{
			ID:       t.ID,
			Name:     t.Value,
			Category: t.Category.Name,
		})
	}
	return result, nil
}

func (ts *TestListService) ListTestLists(ctx context.Context, search string, categoryID *int, page, limit int) (TestListResult, error) {
	page, limit = normalPagination(page, limit)
	items, total, err := ts.repo.ListTestLists(ctx, search, categoryID, page, limit)
	if err != nil {
		return TestListResult{}, err
	}
	totalPages := calculateTotalPage(total, limit)
	return TestListResult{Data: items, Total: total, Page: page, Limit: limit, TotalPages: totalPages}, nil
}

func (ts *TestListService) CreateTestList(ctx context.Context, input models.TestList) (*models.TestList, error) {
	if input.CategoryID < 1 || strings.TrimSpace(input.Value) == "" || strings.TrimSpace(input.Label) == "" {
		return nil, ErrInvalidTestList
	}
	return ts.repo.CreateTestList(ctx, input)
}

func (ts *TestListService) UpdateTestList(ctx context.Context, xid uuid.UUID, input models.TestList) (*models.TestList, error) {
	if input.CategoryID < 1 || strings.TrimSpace(input.Value) == "" || strings.TrimSpace(input.Label) == "" {
		return nil, ErrInvalidTestList
	}
	current, err := ts.repo.GetTestListByXid(ctx, xid)
	if err != nil {
		return nil, err
	}
	isCurrent := true
	input.IsCurrent = &isCurrent
	input.Version = current.Version + 1
	return ts.repo.CreateNewVersion(ctx, xid, input)
}

func (ts *TestListService) DeactivateTestList(ctx context.Context, xid uuid.UUID) (*models.TestList, error) {
	return ts.repo.DeactivateByXid(ctx, xid)
}
