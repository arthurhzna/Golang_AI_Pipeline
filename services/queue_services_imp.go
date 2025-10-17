package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"task_queue/common/aws"
	"task_queue/common/mqtt"
	"task_queue/domain/dto"
	models "task_queue/domain/model"
	"task_queue/repositories"
	"time"
)

type QueueServiceImpl struct {
	awsS3       aws.AWS_S3
	mqtt        mqtt.MQTT
	keyPathAWS  string
	mqttTopic   string
	repository  repositories.QueueRepository
	baseDirSend string
}

func NewQueueService(awsS3 aws.AWS_S3, mqtt mqtt.MQTT, keyPathAWS string, mqttTopic string, repository repositories.QueueRepository, baseDirSend string) QueueService {
	return &QueueServiceImpl{awsS3: awsS3, mqtt: mqtt, keyPathAWS: keyPathAWS, mqttTopic: mqttTopic, repository: repository, baseDirSend: baseDirSend}
}

func (r *QueueServiceImpl) SetQueue(ctx context.Context, data *dto.QueueRequest) (*dto.QueueResponse, error) {

	baseDir := r.baseDirSend
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, err
	}
	timestamp := strings.ReplaceAll(data.Timestamp, "-", "")
	timestamp = strings.ReplaceAll(timestamp, " ", "_")
	timestamp = strings.ReplaceAll(timestamp, ":", "")

	filename := fmt.Sprintf("%s_%s.jpg", data.DeviceId, timestamp)
	filepath := filepath.Join(baseDir, filename)
	src, err := data.Images.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return nil, err
	}

	filepath = strings.ReplaceAll(filepath, "\\", "/")

	DataRedis := models.QueueDataRedis{
		FileName:  filename,
		Path:      filepath,
		DeviceID:  data.DeviceId,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = r.repository.SetQueue(ctx, &DataRedis)
	if err != nil {
		return nil, err
	}

	return &dto.QueueResponse{
		DeviceId:  data.DeviceId,
		Timestamp: data.Timestamp,
	}, nil
}

func (s *QueueServiceImpl) GetQueue(ctx context.Context) (*models.QueuePredictionRedis, error) {
	return s.repository.GetQueue(ctx)
}

func (s *QueueServiceImpl) PublishPredictionToS3AndMQTT(ctx context.Context, data *models.QueuePredictionRedis) error {

	key := fmt.Sprintf("%s/%s", s.keyPathAWS, data.FileName)
	err := s.awsS3.UploadFile(ctx, data.ImageOutputPath, key)
	if err != nil {
		return err
	}

	payload := dto.PredictionMQTTPayload{
		DeviceID:           data.DeviceID,
		Timestamp_In:       data.Timestamp,
		Timestamp_Out:      time.Now().Format("2006-01-02 15:04:05"),
		FileName:           data.FileName,
		ImageAWS3Path:      key,
		OutputText:         data.OutputText,
		PredictedPlatColor: data.PredictedPlatColor,
		PredictedPlatType:  data.PredictedPlatType,
		TimeTakenPredict:   data.TimeTakenPredict,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.mqtt.Publish(ctx, s.mqttTopic, string(payloadJSON))
}
