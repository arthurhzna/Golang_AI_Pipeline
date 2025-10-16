package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWS_S3 interface {
	CreateClient(ctx context.Context) (*s3.Client, error)
	UploadFile(ctx context.Context, s3Client *s3.Client, filePath, key string) error
}
