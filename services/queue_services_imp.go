package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	errWrap "task_queue/common/error"
	"task_queue/domain/dto"
	models "task_queue/domain/model"
	"task_queue/repositories"
)

type QueueServiceImpl struct {
	repository  repositories.QueueRepository
	baseDirSend string
	baseDirGet  string
}

func NewQueueService(repository repositories.QueueRepository, baseDirSend string, baseDirGet string) QueueService {
	return &QueueServiceImpl{repository: repository, baseDirSend: baseDirSend, baseDirGet: baseDirGet}
}

func (r *QueueServiceImpl) SetQueue(ctx context.Context, data *dto.QueueRequest) (*dto.QueueResponse, error) {

	baseDir := r.baseDirSend
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, errWrap.WrapError(fmt.Errorf("failed to create directory: %w", err))
	}
	timestamp := strings.ReplaceAll(data.Timestamp, "-", "")
	timestamp = strings.ReplaceAll(timestamp, " ", "_")
	timestamp = strings.ReplaceAll(timestamp, ":", "")

	filename := fmt.Sprintf("%s_%s.jpg", data.DeviceId, timestamp)
	filepath := filepath.Join(baseDir, filename)
	src, err := data.Images.Open()
	if err != nil {
		return nil, errWrap.WrapError(fmt.Errorf("failed to open uploaded file: %w", err))
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return nil, errWrap.WrapError(fmt.Errorf("failed to create file: %w", err))
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return nil, errWrap.WrapError(fmt.Errorf("failed to save file: %w", err))
	}

	DataRedis := models.QueueDataRedis{
		Path:      filepath,
		DeviceID:  data.DeviceId,
		Timestamp: data.Timestamp,
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
