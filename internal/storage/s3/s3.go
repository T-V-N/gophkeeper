package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Store struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	Cfg           *config.Config
	BucketName    string
}

func InitS3Storage(ctx context.Context, cfg *config.Config) *S3Store {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf(cfg.S3URL),
		}, nil
	})

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3Secret, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: &cfg.S3Bucket})

	if err != nil {
		var bne *types.BucketAlreadyExists
		if errors.As(err, &bne) {
			log.Println("error:", bne)
		} else {
			log.Panic("unable to connect to the s3")
		}
	}

	return &S3Store{S3Client: client, BucketName: cfg.S3Bucket, Cfg: cfg}
}

func (s S3Store) UploadLargeObject(ctx context.Context, objectKey string, largeObject []byte) error {
	var partMiBs int64 = 10

	largeBuffer := bytes.NewReader(largeObject)

	uploader := manager.NewUploader(s.S3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
		Body:   largeBuffer,
	})

	return err
}

func (s S3Store) DeleteObject(ctx context.Context, objectKey string) error {
	_, err := s.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{Key: aws.String(objectKey)})

	return err
}

func (s S3Store) GetUploadLink(ctx context.Context, objectKey string) (string, error) {
	request, err := s.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(int64(s.Cfg.FileUpdateTimeWindow) * int64(time.Second))
	})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}
