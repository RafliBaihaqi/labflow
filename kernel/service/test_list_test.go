package service

import (
	"context"
	"errors"
	"testing"

	"labflow/apperr"
	"labflow/models"

	"github.com/google/uuid"
)

type fakeTestListRepo struct {
	current *models.TestList
	got     *models.TestList
}

func (f *fakeTestListRepo) GetTestListByCategory(ctx context.Context, categoryID int) ([]models.TestList, error) {
	return nil, nil
}

func (f *fakeTestListRepo) ListTestLists(ctx context.Context, search string, categoryID *int, page, limit int) ([]models.TestList, int64, error) {
	return nil, 0, nil
}

func (f *fakeTestListRepo) CreateTestList(ctx context.Context, input models.TestList) (*models.TestList, error) {
	return &input, nil
}

func (f *fakeTestListRepo) GetTestListByXid(ctx context.Context, xid uuid.UUID) (*models.TestList, error) {
	if f.current == nil {
		return nil, apperr.ErrNotFound
	}
	return f.current, nil
}

func (f *fakeTestListRepo) CreateNewVersion(ctx context.Context, xid uuid.UUID, input models.TestList) (*models.TestList, error) {
	f.got = &input
	return &input, nil
}

func (f *fakeTestListRepo) DeactivateByXid(ctx context.Context, xid uuid.UUID) (*models.TestList, error) {
	return f.current, nil
}

func TestUpdateTestListBumpsVersion(t *testing.T) {
	v := 3
	repo := &fakeTestListRepo{current: &models.TestList{Version: v}}
	svc := NewTestListService(repo)

	_, err := svc.UpdateTestList(context.Background(), uuid.New(), models.TestList{
		CategoryID: 1,
		Value:      "ph",
		Label:      "pH",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repo.got == nil {
		t.Fatal("CreateNewVersion was not called")
	}
	if repo.got.Version != v+1 {
		t.Fatalf("Version = %d, want %d", repo.got.Version, v+1)
	}
	if repo.got.IsCurrent == nil || !*repo.got.IsCurrent {
		t.Fatal("new version must be current")
	}
}

func TestUpdateTestListRejectsInvalidInput(t *testing.T) {
	svc := NewTestListService(&fakeTestListRepo{current: &models.TestList{Version: 1}})

	if _, err := svc.UpdateTestList(context.Background(), uuid.New(), models.TestList{Value: "ph"}); !errors.Is(err, ErrInvalidTestList) {
		t.Fatalf("got %v, want ErrInvalidTestList", err)
	}
}
