package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"task_queue/common/aws"
	errWrap "task_queue/common/error"
	"task_queue/common/mqtt"
	"task_queue/domain/dto"
	"task_queue/repositories"
	"time"
)

type WorkerImpl struct {
	repository   repositories.QueueRepository
	awsS3        aws.AWS_S3
	key_path_aws string
	mqtt         mqtt.MQTT
	mqtt_topic   string
}

func NewWorker(repository repositories.QueueRepository, awsS3 aws.AWS_S3, key_aws_bucket string, mqtt mqtt.MQTT, mqtt_topic string) Worker {
	return &WorkerImpl{repository: repository, awsS3: awsS3, key_path_aws: key_aws_bucket, mqtt: mqtt, mqtt_topic: mqtt_topic}
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
						time.Sleep(50 * time.Millisecond)
						continue
					}
					if data == nil {
						time.Sleep(50 * time.Millisecond)
						continue
					}

					key := fmt.Sprintf("%s/%s", w.key_path_aws, data.FileName)
					err = w.awsS3.UploadFile(ctx, data.ImageOutputPath, key)
					if err != nil {
						// time.Sleep(1 * time.Second)
						continue
					}
					fmt.Printf("Worker %d - file uploaded to AWS S3\n", workerID)

					payload := dto.PredictionMQTTPayload{
						DeviceID:           data.DeviceID,
						Timestamp_In:       data.Timestamp,
						Timestamp_Out:      time.Now().Format("2006-01-02 15:04:05"),
						FileName:           data.FileName,
						ImageOutputPath:    key,
						OutputText:         data.OutputText,
						PredictedPlatColor: data.PredictedPlatColor,
						PredictedPlatType:  data.PredictedPlatType,
						TimeTakenPredict:   data.TimeTakenPredict,
					}

					payloadJSON, err := json.Marshal(payload)
					if err != nil {
						errWrap.WrapError(fmt.Errorf("failed to marshal payload to JSON: %w", err))
						continue
					}

					err = w.mqtt.Publish(ctx, w.mqtt_topic, string(payloadJSON))
					if err != nil {
						continue
					}
					fmt.Printf("Worker %d - payload: %s\n", workerID, string(payloadJSON))
					time.Sleep(30 * time.Millisecond)
				}
			}
		}(i + 1)
	}
	return nil
}
