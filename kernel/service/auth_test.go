package service

import (
	"context"
	"labflow/auth"
	"labflow/models"
	"testing"
)

type fakeAuthRepo struct {
	authRepo
	updatedPassword string
	updateErr       error
}

func (f *fakeAuthRepo) UpdateUserPassword(ctx context.Context, email, hashedPassword string) (*models.User, error) {
	f.updatedPassword = hashedPassword
	return &models.User{Email: email}, f.updateErr
}

type fakeTokenManager struct {
	tokenManager
	parseEmail string
	parseErr   error
}

func (f *fakeTokenManager) ParseResetToken(ctx context.Context, tokenString string) (string, error) {
	return f.parseEmail, f.parseErr
}

func TestResetPasswordHashesPassword(t *testing.T) {
	repo := &fakeAuthRepo{}
	tokenUtils := &fakeTokenManager{parseEmail: "test@example.com"}
	svc := NewService(repo, nil, nil, tokenUtils, 1)

	_, err := svc.ResetPassword(context.Background(), "valid-token", "myPlaintextPassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.updatedPassword == "" {
		t.Fatal("expected password to be updated in repo, but it was empty")
	}

	if repo.updatedPassword == "myPlaintextPassword" {
		t.Fatal("expected password to be hashed in repo, but got plaintext")
	}

	err = auth.ComparePassword(repo.updatedPassword, "myPlaintextPassword")
	if err != nil {
		t.Fatalf("expected repo.updatedPassword to match 'myPlaintextPassword', error: %v", err)
	}
}
