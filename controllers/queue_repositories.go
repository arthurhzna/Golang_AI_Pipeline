package services

import (
	"context"
	"task_queue/domain/dto"
)

type QueueService interface {
	SetQueue(context.Context, *dto.QueueRequest) error
}
