package s3

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type S3Store struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	Cfg           *config.S3Config
	BucketName    string
}

func InitS3Storage(ctx context.Context, cfg *config.S3Config) *S3Store {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           fmt.Sprintf(cfg.S3URL),
			SigningRegion: region,
		}, nil
	})

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3Secret, "")),
		awsConfig.WithRegion("eu-central-1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

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

	return utils.WrapError(err, utils.ErrThirdParty)
}

func (s S3Store) DeleteObject(ctx context.Context, objectKey string) error {
	_, err := s.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{Key: aws.String(objectKey)})

	return utils.WrapError(err, utils.ErrThirdParty)
}

func (s S3Store) GetUploadLink(ctx context.Context, objectKey string) (string, error) {
	request, err := s.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(int64(s.Cfg.FileUpdateTimeWindow) * int64(time.Second))
	})

	if err != nil {
		return "", utils.WrapError(err, utils.ErrThirdParty)
	}

	return request.URL, nil
}

func (s S3Store) GetFileUpdatedAt(ctx context.Context, objectKey string) (time.Time, error) {
	request, err := s.S3Client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return time.Time{}, utils.WrapError(err, utils.ErrThirdParty)
	}

	return *request.LastModified, nil
}
