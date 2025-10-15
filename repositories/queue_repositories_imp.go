package repositories

import (
	"context"
	"encoding/json"
	"strings"
	errWrap "task_queue/common/error"
	errConstant "task_queue/constants/error"
	models "task_queue/domain/model"

	"github.com/go-redis/redis/v8"
)

type QueueRepositoryImpl struct {
	db                *redis.Client
	keyRedisGroupSend string
}

func NewQueueRepository(db *redis.Client, key_redis_group_send string) QueueRepository {
	return &QueueRepositoryImpl{db: db, keyRedisGroupSend: key_redis_group_send}
}

func (r *QueueRepositoryImpl) SetQueue(ctx context.Context, data *models.QueueDataRedis) error {
	data.Path = strings.ReplaceAll(data.Path, "\\", "/")
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
