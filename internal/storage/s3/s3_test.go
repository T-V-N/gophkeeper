package s3_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"

	storage "github.com/T-V-N/gophkeeper/internal/storage/s3"
)

func InitTestConfig() config.Config {
	return config.Config{S3URL: "http://localhost:9090", S3AccessKey: "hey", S3Secret: "hey2", S3Bucket: "mockbucket"}
}
func Test_ConnectToS3(t *testing.T) {
	cfg := InitTestConfig()

	t.Run("Connects to S3", func(t *testing.T) {
		S3 := storage.InitS3Storage(context.Background(), &cfg)
		buckets, err := S3.S3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

		bucketList := []string{}
		for _, b := range buckets.Buckets {
			bucketList = append(bucketList, *b.Name)
		}
		assert.NoError(t, err, "Shall connect and list buckets")
		assert.Contains(t, bucketList, "test", "Shall containt test bucket")
	})
}

func Test_S3FileUpload(t *testing.T) {
	cfg := InitTestConfig()

	t.Run("Connects to S3", func(t *testing.T) {
		file := make([]byte, 1024)
		rand.Read(file)

		S3 := storage.InitS3Storage(context.Background(), &cfg)

		err := S3.UploadLargeObject(context.Background(), "test-file.txt", file)
		fmt.Print(err)
	})
}
