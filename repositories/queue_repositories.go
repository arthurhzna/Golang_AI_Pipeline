package repositories

import (
	"context"
	"time"
)

type QueueRepository interface {
	SetQueue(context.Context, string, string, time.Time) error
}
