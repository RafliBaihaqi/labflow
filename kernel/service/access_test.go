package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"labflow/models"
)

type fakeAccessRepo struct {
	calls int
	err   error
	perms []string
}

func (f *fakeAccessRepo) GetUserByID(_ context.Context, _ int) (*models.User, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	var accesses []models.Access
	for i := range f.perms {
		accesses = append(accesses, models.Access{ID: i + 1, AccessName: &f.perms[i]})
	}
	return &models.User{Roles: []models.Role{{ID: 1, Accesses: accesses}}}, nil
}

func TestHasAccessCachesPermissions(t *testing.T) {
	repo := &fakeAccessRepo{perms: []string{"read_reservation"}}
	svc := NewAccessServiceWithTTL(repo, time.Minute)

	for i := 0; i < 5; i++ {
		ok, err := svc.HasAccess(context.Background(), 7, "read_reservation")
		if err != nil || !ok {
			t.Fatalf("expected access, got ok=%v err=%v", ok, err)
		}
	}
	if ok, _ := svc.HasAccess(context.Background(), 7, "delete_test"); ok {
		t.Fatal("expected delete_test to be denied")
	}
	if repo.calls != 1 {
		t.Fatalf("expected 1 repo call for 6 checks, got %d", repo.calls)
	}
}

func TestHasAccessReloadsAfterTTL(t *testing.T) {
	repo := &fakeAccessRepo{perms: []string{"read_test"}}
	svc := NewAccessServiceWithTTL(repo, time.Second)
	now := time.Unix(1000, 0)
	svc.now = func() time.Time { return now }

	_, _ = svc.HasAccess(context.Background(), 1, "read_test")
	now = now.Add(2 * time.Second)
	_, _ = svc.HasAccess(context.Background(), 1, "read_test")

	if repo.calls != 2 {
		t.Fatalf("expected reload after TTL, got %d calls", repo.calls)
	}
}

func TestHasAccessInvalidate(t *testing.T) {
	repo := &fakeAccessRepo{perms: []string{"read_test"}}
	svc := NewAccessServiceWithTTL(repo, time.Hour)

	_, _ = svc.HasAccess(context.Background(), 1, "read_test")
	svc.Invalidate(1)
	_, _ = svc.HasAccess(context.Background(), 1, "read_test")

	if repo.calls != 2 {
		t.Fatalf("expected reload after Invalidate, got %d calls", repo.calls)
	}
}

func TestHasAccessErrorIsNotADenial(t *testing.T) {
	repo := &fakeAccessRepo{err: errors.New("too many clients")}
	svc := NewAccessServiceWithTTL(repo, time.Minute)

	ok, err := svc.HasAccess(context.Background(), 1, "read_test")
	if err == nil || ok {
		t.Fatalf("expected an error and no access, got ok=%v err=%v", ok, err)
	}
	// failures must not be cached
	repo.err = nil
	repo.perms = []string{"read_test"}
	if ok, err := svc.HasAccess(context.Background(), 1, "read_test"); err != nil || !ok {
		t.Fatalf("expected recovery after DB error, got ok=%v err=%v", ok, err)
	}
}
