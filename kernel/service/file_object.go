package service

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

var ErrObjectNotFound = errors.New("object not found")

type UploadDestination string

const (
	DestinationReservation UploadDestination = "reservation"
	DestinationCertificate UploadDestination = "certificate"
	DestinationUser        UploadDestination = "user"
)

type objectGetter interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type BucketBasics struct {
	S3Client objectGetter
}

func NewFileObjectService(client objectGetter) *BucketBasics {
	return &BucketBasics{S3Client: client}
}

type ObjectData struct {
	Body        io.ReadCloser
	ContentType string
	Size        int64
}

func (b BucketBasics) GetObject(ctx context.Context, bucketName, objectKey string) (*ObjectData, error) {
	out, err := b.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && (apiErr.ErrorCode() == "NoSuchKey" || apiErr.ErrorCode() == "NotFound") {
			return nil, ErrObjectNotFound
		}
		log.Printf("couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return nil, err
	}
	return &ObjectData{
		Body:        out.Body,
		ContentType: aws.ToString(out.ContentType),
		Size:        aws.ToInt64(out.ContentLength),
	}, nil
}
