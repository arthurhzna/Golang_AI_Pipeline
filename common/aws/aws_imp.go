package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWS_S3_Impl struct {
	accessKeyID     string
	secretAccessKey string
	region          string
	bucket          string
}

func NewAWS_S3(accessKeyID string, secretAccessKey string, defaultRegion string, bucket string) AWS_S3 {
	return &AWS_S3_Impl{accessKeyID: accessKeyID, secretAccessKey: secretAccessKey, region: defaultRegion, bucket: bucket}
}

func (a *AWS_S3_Impl) CreateClient(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(a.region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				a.accessKeyID,
				a.secretAccessKey,
				"")))
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

func (a *AWS_S3_Impl) UploadFile(ctx context.Context, s3Client *s3.Client, filePath, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("image/jpeg"),
	})
	return err
}
