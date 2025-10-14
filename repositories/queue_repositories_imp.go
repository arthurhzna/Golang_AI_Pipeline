package repositories

import (
	"context"
	"encoding/json"

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
	dataRedis, err := json.Marshal(data)
	if err != nil {
		return errWrap.WrapError(errConstant.ErrInternalServerError)
	}
	err = r.db.LPush(ctx, "queue_image", string(dataRedis)).Err()
	if err != nil {
		return errWrap.WrapError(errConstant.ErrInternalServerError)
	}
	return nil
}
