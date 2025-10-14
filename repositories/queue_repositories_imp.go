package repositories

import (
	"context"
	errWrap "task_queue/common/error"
	errConstant "task_queue/constants/error"
	models "task_queue/domain/model"

	"github.com/go-redis/redis/v8"
)

type QueueRepositoryImpl struct {
	db *redis.Client
}

func NewQueueRepository(db *redis.Client) QueueRepository {
	return &QueueRepositoryImpl{db: db}
}

func (r *QueueRepositoryImpl) SetQueue(ctx context.Context, data *models.QueueDataRedis) error {
	// dataRedis := models.QueueDataRedis{
	// 	Path:      data.Path,
	// 	DeviceID:  data.DeviceID,
	// 	Timestamp: data.Timestamp,
	// }
	err := r.db.Set(ctx, "queue_image", data, 0).Err()
	if err != nil {
		return errWrap.WrapError(errConstant.ErrInternalServerError)
	}
	return nil
}
