package service

import (
	"context"
	"errors"
	"testing"

	"labflow/models"

	"github.com/google/uuid"
)

type fakeTestRepo struct {
	test    *models.Test
	getErr  error
	updated bool
}

func (f *fakeTestRepo) GetTestByXid(ctx context.Context, xid uuid.UUID) (*models.Test, error) {
	return f.test, f.getErr
}

func (f *fakeTestRepo) UpdateTestByXid(ctx context.Context, xid uuid.UUID, testResult, testSummary, notes *string, updatedBy int) (*models.Test, error) {
	f.updated = true
	return f.test, nil
}

func (f *fakeTestRepo) UpdateIsApproved(ctx context.Context, xid uuid.UUID, isApproved *bool, updatedBy int) (*models.Test, error) {
	return f.test, nil
}

type fakeRoleReader struct {
	user *models.User
	err  error
}

func (f *fakeRoleReader) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	return f.user, f.err
}

func testWithCategory(category string) *models.Test {
	return &models.Test{TestList: models.TestList{Category: models.Category{Name: category}}}
}

func userWithRoles(names ...string) *models.User {
	roles := make([]models.Role, 0, len(names))
	for _, n := range names {
		name := n
		roles = append(roles, models.Role{RoleName: &name})
	}
	return &models.User{Roles: roles}
}

func TestUpdateTestByXidRoleAccess(t *testing.T) {
	cases := []struct {
		name     string
		roles    []string
		category string
		wantErr  error
	}{
		{"admin passes any category", []string{"admin"}, "Chemistry", nil},
		{"micro-analyst passes microbiology", []string{"micro-analyst"}, "Microbiology", nil},
		{"chemical-qc passes chemistry", []string{"chemical-qc"}, "Chemistry", nil},
		{"micro-analyst rejected on chemistry", []string{"micro-analyst"}, "Chemistry", ErrUnauthorizedRole},
		{"customer rejected", []string{"customer"}, "Microbiology", ErrUnauthorizedRole},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeTestRepo{test: testWithCategory(tc.category)}
			svc := NewTestService(repo, &fakeRoleReader{user: userWithRoles(tc.roles...)})

			_, err := svc.UpdateTestByXid(context.Background(), uuid.New(), 1, UpdateTestParams{})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got err %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && !repo.updated {
				t.Fatal("expected repo update to be called")
			}
		})
	}
}
