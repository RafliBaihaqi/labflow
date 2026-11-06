package service

import (
	"context"
	"errors"
	"labflow/bootstrap"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type fakePutter struct {
	got *s3.PutObjectInput
	err error
}

func (f *fakePutter) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	f.got = params
	return &s3.PutObjectOutput{}, f.err
}

func TestUploadServiceUpload(t *testing.T) {
	f := &fakePutter{}
	svc := NewUploadService(f, bootstrap.BucketReservations, bootstrap.BucketCertificates, bootstrap.BucketUsers)

	name, err := svc.Upload(context.Background(), DestinationCertificate, strings.NewReader("PDF"), 3, "application/pdf", "pdf")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if f.got == nil || aws.ToString(f.got.Bucket) != bootstrap.BucketCertificates {
		t.Fatalf("expected certificate bucket, got %+v", f.got)
	}
	if !strings.HasSuffix(name, ".pdf") || name != aws.ToString(f.got.Key) {
		t.Fatalf("unexpected key: %q", name)
	}
	if aws.ToInt64(f.got.ContentLength) != 3 || aws.ToString(f.got.ContentType) != "application/pdf" {
		t.Fatalf("unexpected length/content-type: %+v", f.got)
	}
}

func TestUploadServiceUnknownDestination(t *testing.T) {
	svc := NewUploadService(&fakePutter{}, "r", "c", "u")
	if _, err := svc.Upload(context.Background(), "nope", strings.NewReader("x"), 1, "text/plain", "txt"); !errors.Is(err, ErrUnknownDestination) {
		t.Fatalf("expected ErrUnknownDestination, got %v", err)
	}
}

func TestUploadServicePutError(t *testing.T) {
	svc := NewUploadService(&fakePutter{err: errors.New("boom")}, "r", "c", "u")
	if _, err := svc.Upload(context.Background(), DestinationUser, strings.NewReader("x"), 1, "image/png", "png"); !errors.Is(err, ErrUploadFailed) {
		t.Fatalf("expected ErrUploadFailed, got %v", err)
	}
}
