package workers

import (
	"context"
	"fmt"
	"strings"
	"task_queue/repositories"

	"task_queue/common/aws"
	errWrap "task_queue/common/error"
)

type WorkerImpl struct {
	repository repositories.QueueRepository
	awsS3      aws.AWS_S3
}

func NewWorker(repository repositories.QueueRepository, awsS3 aws.AWS_S3) Worker {
	return &WorkerImpl{repository: repository, awsS3: awsS3}
}

func (w *WorkerImpl) Run(ctx context.Context, worker int) error {
	s3Client, err := w.awsS3.CreateClient(ctx)
	if err != nil {
		return errWrap.WrapError(fmt.Errorf("failed to create AWS S3 client: %w", err))
	}
	for i := 0; i < worker; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					errWrap.WrapError(fmt.Errorf("worker %d stopped", workerID))
					return
				default:
					data, err := w.repository.GetQueue(ctx)
					if err != nil {
						errWrap.WrapError(err)
						// time.Sleep(1 * time.Second)
						continue
					}
					if data == nil {
						// time.Sleep(1 * time.Second)
						continue
					}
					fmt.Printf("Worker %d - data: %+v\n", workerID, data)
					data.OutputPath = strings.ReplaceAll(data.OutputPath, "\\", "/")

					err = w.awsS3.UploadFile(ctx, s3Client, data.OutputPath, data.OutputPath)
					if err != nil {
						errWrap.WrapError(fmt.Errorf("failed to upload file to AWS S3: %w", err))
						// time.Sleep(1 * time.Second)
						continue
					}
				}
			}
		}(i + 1)
	}
	return nil
}
