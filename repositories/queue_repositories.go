package repositories

import (
	"context"
	models "task_queue/domain/model"
)

type QueueRepository interface {
	SetQueue(context.Context, *models.QueueDataRedis) error
}
