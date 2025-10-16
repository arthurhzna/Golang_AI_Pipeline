package aws

import (
	"context"
)

type AWS_S3 interface {
	CreateClient(ctx context.Context) error
	UploadFile(ctx context.Context, filePath, key string) error
}
