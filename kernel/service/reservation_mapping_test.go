package service

import (
	"context"
	"testing"

	"labflow/models"
)

func TestToReservationResultMapsFileURLs(t *testing.T) {
	rs := &ReservationService{}

	sample := "sample.png"
	avatar := "avatar.jpg"
	r := models.Reservation{
		SamplePhoto: &sample,
		User:        models.User{Avatar: &avatar},
	}

	got := rs.toReservationResult(context.Background(), r)

	if got.SamplePhoto == nil || *got.SamplePhoto != "/api/files/reservations/sample.png" {
		t.Fatalf("SamplePhoto = %v, want /api/files/reservations/sample.png", got.SamplePhoto)
	}
	if got.User.Avatar == nil || *got.User.Avatar != "/api/files/users/avatar.jpg" {
		t.Fatalf("User.Avatar = %v, want /api/files/users/avatar.jpg", got.User.Avatar)
	}
}
