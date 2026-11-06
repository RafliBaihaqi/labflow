package service

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type fakeGetter struct {
	out *s3.GetObjectOutput
	err error
}

func (f *fakeGetter) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return f.out, f.err
}

func TestFileObjectServiceGetObjectSuccess(t *testing.T) {
	f := &fakeGetter{out: &s3.GetObjectOutput{
		Body:          io.NopCloser(strings.NewReader("PNGBYTES")),
		ContentType:   aws.String("image/png"),
		ContentLength: aws.Int64(8),
	}}
	svc := NewFileObjectService(f)

	data, err := svc.GetObject(context.Background(), "reservations", "abc.png")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if data.ContentType != "image/png" || data.Size != 8 {
		t.Fatalf("unexpected data: %+v", data)
	}
}

func TestFileObjectServiceGetObjectNotFound(t *testing.T) {
	f := &fakeGetter{err: &types.NoSuchKey{}}
	svc := NewFileObjectService(f)

	_, err := svc.GetObject(context.Background(), "reservations", "missing.png")
	if !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("expected ErrObjectNotFound, got %v", err)
	}
}

func TestFileObjectServiceGetObjectOtherError(t *testing.T) {
	f := &fakeGetter{err: errors.New("connection refused")}
	svc := NewFileObjectService(f)

	_, err := svc.GetObject(context.Background(), "reservations", "x.png")
	if err == nil || errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("expected a non-not-found error, got %v", err)
	}
}
