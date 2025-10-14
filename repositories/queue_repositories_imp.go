package repositories

import (
	"context"
	errConstant "task_queue/constants/error"
	"task_queue/domain/model"
	"time"

	"github.com/go-redis/redis/v8"
)

type QueueRepositoryImpl struct {
	db *redis.Client
}

func NewQueueRepository(db *redis.Client) QueueRepository {
	return &QueueRepositoryImpl{db: db}
}

func (r *QueueRepositoryImpl) SetQueue(ctx context.Context, path string, device_id string, timestamp time.Time) error {
	data := models.QueueDataRedis{
		Path:      path,
		DeviceID:  device_id,
		Timestamp: timestamp,
	}
	err := r.db.Set(ctx, "queue_image", data, 0).Err()
	if err != nil {
		return errConstant.ErrInternalServerError
	}
	return nil
}
