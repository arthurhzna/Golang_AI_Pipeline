package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	errWrap "task_queue/common/error"
)

type AWS_S3_Impl struct {
	client          *s3.Client
	accessKeyID     string
	secretAccessKey string
	region          string
	bucket          string
}

func NewAWS_S3(accessKeyID string, secretAccessKey string, defaultRegion string, bucket string) AWS_S3 {
	return &AWS_S3_Impl{accessKeyID: accessKeyID, secretAccessKey: secretAccessKey, region: defaultRegion, bucket: bucket}
}

func (a *AWS_S3_Impl) CreateClient(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(a.region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				a.accessKeyID,
				a.secretAccessKey,
				"")))
	if err != nil {
		return errWrap.WrapError(err)
	}
	a.client = s3.NewFromConfig(cfg)
	return nil
}

func (a *AWS_S3_Impl) UploadFile(ctx context.Context, filePath, key string) error {
	if a.client == nil {
		return errWrap.WrapError(fmt.Errorf("S3 client is not initialized"))
	}
	file, err := os.Open(filePath)
	if err != nil {
		return errWrap.WrapError(err)
	}
	defer file.Close()

	_, err = a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("image/jpeg"),
	})

	return err
}
