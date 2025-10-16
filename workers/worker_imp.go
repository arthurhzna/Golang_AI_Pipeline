package workers

import (
	"context"
	"fmt"
	"strings"
	"task_queue/repositories"

	"task_queue/common/aws"
	errWrap "task_queue/common/error"
	"task_queue/common/mqtt"
)

type WorkerImpl struct {
	repository repositories.QueueRepository
	awsS3      aws.AWS_S3
	mqtt       mqtt.MQTT
}

func NewWorker(repository repositories.QueueRepository, awsS3 aws.AWS_S3, key_aws_bucket string, mqtt mqtt.MQTT) Worker {
	return &WorkerImpl{repository: repository, awsS3: awsS3, mqtt: mqtt}
}

func (w *WorkerImpl) Run(ctx context.Context, worker int) error {

	err_aws := w.awsS3.CreateClient(ctx)
	if err_aws != nil {
		return err_aws
	}
	err_mqtt := w.mqtt.Connect(ctx)
	if err_mqtt != nil {
		return err_mqtt
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
						// time.Sleep(1 * time.Second)
						continue
					}
					if data == nil {
						// time.Sleep(1 * time.Second)
						continue
					}
					fmt.Printf("Worker %d - data: %+v\n", workerID, data)
					data.OutputPath = strings.ReplaceAll(data.OutputPath, "\\", "/")

					err = w.awsS3.UploadFile(ctx, data.OutputPath, data.OutputPath)
					if err != nil {
						errWrap.WrapError(fmt.Errorf("failed to upload file to AWS S3: %w", err))
						// time.Sleep(1 * time.Second)
						continue
					}
					fmt.Printf("Worker %d - file uploaded to AWS S3\n", workerID)
				}
			}
		}(i + 1)
	}
	return nil
}
