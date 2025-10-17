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
	db                *redis.Client
	keyRedisGroupSend string
	keyRedisGroupGet  string
}

func NewQueueRepository(db *redis.Client, key_redis_group_send string, key_redis_group_get string) QueueRepository {
	return &QueueRepositoryImpl{db: db, keyRedisGroupSend: key_redis_group_send, keyRedisGroupGet: key_redis_group_get}
}

func (r *QueueRepositoryImpl) SetQueue(ctx context.Context, data *models.QueueDataRedis) error {
	dataRedis, err := json.Marshal(data)
	if err != nil {
		return errWrap.WrapError(errConstant.ErrInternalServerError)
	}
	err = r.db.LPush(ctx, r.keyRedisGroupSend, string(dataRedis)).Err()
	if err != nil {
		return errWrap.WrapError(errConstant.ErrInternalServerError)
	}
	return nil
}

func (r *QueueRepositoryImpl) GetQueue(ctx context.Context) (*models.QueuePredictionRedis, error) {
	data, err := r.db.RPop(ctx, r.keyRedisGroupGet).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, errWrap.WrapError(errConstant.ErrQueueNotFound)
	}
	var prediction models.QueuePredictionRedis
	err = json.Unmarshal([]byte(data), &prediction)
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrInternalServerError)
	}
	return &prediction, nil
}
