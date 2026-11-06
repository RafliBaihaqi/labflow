package service

import "testing"

func strPtr(s string) *string { return &s }

func TestFileURL(t *testing.T) {
	t.Run("nil object returns nil", func(t *testing.T) {
		if got := FileURL(FileCategoryReservation, nil); got != nil {
			t.Fatalf("expected nil, got %v", *got)
		}
	})

	t.Run("empty object returns same pointer unchanged", func(t *testing.T) {
		empty := strPtr("")
		got := FileURL(FileCategoryUser, empty)
		if got == nil || *got != "" {
			t.Fatalf("expected empty string pointer, got %v", got)
		}
	})

	t.Run("reservation object builds relative path", func(t *testing.T) {
		got := FileURL(FileCategoryReservation, strPtr("abc.png"))
		if got == nil || *got != "/api/files/reservations/abc.png" {
			t.Fatalf("unexpected url: %v", got)
		}
	})

	t.Run("user object builds relative path", func(t *testing.T) {
		got := FileURL(FileCategoryUser, strPtr("avatar.jpg"))
		if got == nil || *got != "/api/files/users/avatar.jpg" {
			t.Fatalf("unexpected url: %v", got)
		}
	})

	t.Run("certificate object builds relative path", func(t *testing.T) {
		got := FileURL(FileCategoryCertificate, strPtr("cert.pdf"))
		if got == nil || *got != "/api/files/certificates/cert.pdf" {
			t.Fatalf("unexpected url: %v", got)
		}
	})
}
