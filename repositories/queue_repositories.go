package repositories

import (
	"context"
	"task_queue/domain/model"
)

type QueueRepository interface {
	SetQueue(context.Context, *models.QueueDataRedis) error
}
