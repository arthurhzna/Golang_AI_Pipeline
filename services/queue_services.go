package services

import (
	"context"
	"task_queue/domain/dto"
	models "task_queue/domain/model"
)

type QueueService interface {
	SetQueue(context.Context, *dto.QueueRequest) (*dto.QueueResponse, error)
	GetQueue(context.Context) (*models.QueuePredictionRedis, error)
	PublishPredictionToS3AndMQTT(context.Context, *models.QueuePredictionRedis) error
}
