package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

const (
	BucketReservations = "reservations"
	BucketUsers        = "users"
	BucketCertificates = "certificates"
)

// storageConfig resolves the S3-compatible object store from the environment.
// Precedence: S3_ENDPOINT (any S3-compatible store), then MINIO_ENDPOINT (the
// docker-compose stacks), then Cloudflare R2 via ACCOUNT_ID.
type storageConfig struct {
	endpoint  string
	accessKey string
	secretKey string
	pathStyle bool
}

func resolveStorageConfig() storageConfig {
	if ep := os.Getenv("S3_ENDPOINT"); ep != "" {
		return storageConfig{endpoint: ep, accessKey: os.Getenv("ACCESS_KEY_ID"), secretKey: os.Getenv("SECRET_ACCESS_KEY"), pathStyle: true}
	}
	if ep := os.Getenv("MINIO_ENDPOINT"); ep != "" {
		if !strings.HasPrefix(ep, "http://") && !strings.HasPrefix(ep, "https://") {
			ep = "http://" + ep
		}
		return storageConfig{endpoint: ep, accessKey: os.Getenv("MINIO_ACCESS"), secretKey: os.Getenv("MINIO_SECRET"), pathStyle: true}
	}
	return storageConfig{
		endpoint:  fmt.Sprintf("https://%s.r2.cloudflarestorage.com", os.Getenv("ACCOUNT_ID")),
		accessKey: os.Getenv("ACCESS_KEY_ID"),
		secretKey: os.Getenv("SECRET_ACCESS_KEY"),
	}
}

func ConnectToR2(ctx context.Context) *s3.Client {
	sc := resolveStorageConfig()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(sc.accessKey, sc.secretKey, "")),
		config.WithRegion("auto"), // Required by SDK but not used by R2/MinIO
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(sc.endpoint)
		o.UsePathStyle = sc.pathStyle
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})
	return client
}

type BucketBasics struct {
	S3Client *s3.Client
}

func (b BucketBasics) CreateBucket(ctx context.Context, name string, region string) error {
	_, err := b.S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			log.Printf("You already own bucket %s.\n", name)
			err = owned
		} else if errors.As(err, &exists) {
			log.Printf("Bucket %s already exists.\n", name)
			err = exists
		}
	} else {
		err = s3.NewBucketExistsWaiter(b.S3Client).Wait(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(name)}, time.Minute)
		if err != nil {
			log.Printf("Failed to attempt to wait for bucket %s to exist.\n", name)
		}
	}
	return err
}

func (b BucketBasics) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := b.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Printf("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}

func InitBuckets(ctx context.Context, s3Client *s3.Client) {
	bucketBasics := BucketBasics{S3Client: s3Client}

	for _, name := range []string{BucketReservations, BucketUsers, BucketCertificates} {
		bucketExists, err := bucketBasics.BucketExists(ctx, name)
		if err != nil {
			log.Fatalf("Failed to check if %s bucket exists: %v", name, err)
		}
		if !bucketExists {
			if err := bucketBasics.CreateBucket(ctx, name, ""); err != nil {
				log.Fatalf("Failed to create %s bucket: %v", name, err)
			}
		}
		log.Println("Buckets initialized successfully.")
	}
	log.Println("Bucket init completed")
}
