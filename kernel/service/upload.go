package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var (
	ErrUploadFailed       = errors.New("failed to upload file")
	ErrUnknownDestination = errors.New("unknown upload destination")
)

type objectPutter interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type UploadService struct {
	client            objectPutter
	bucketReservation string
	bucketCertificate string
	bucketUser        string
}

func NewUploadService(client objectPutter, bucketReservation, bucketCertificate, bucketUser string) *UploadService {
	return &UploadService{
		client:            client,
		bucketReservation: bucketReservation,
		bucketCertificate: bucketCertificate,
		bucketUser:        bucketUser,
	}
}

func (s *UploadService) Upload(ctx context.Context, dest UploadDestination, file io.Reader, size int64, contentType, extension string) (string, error) {
	bucket, ok := s.bucketFor(dest)
	if !ok {
		return "", ErrUnknownDestination
	}
	objectName := fmt.Sprintf("%s.%s", uuid.New().String(), extension)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objectName),
		Body:          file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		log.Printf("couldn't upload object %v:%v. Here's why: %v\n", bucket, objectName, err)
		return "", ErrUploadFailed
	}
	return objectName, nil
}

func (s *UploadService) bucketFor(dest UploadDestination) (string, bool) {
	switch dest {
	case DestinationReservation:
		return s.bucketReservation, true
	case DestinationCertificate:
		return s.bucketCertificate, true
	case DestinationUser:
		return s.bucketUser, true
	default:
		return "", false
	}
}
