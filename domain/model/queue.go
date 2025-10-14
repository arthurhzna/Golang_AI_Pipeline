package models

import (
	"time"
)

type QueueDataRedis struct {
	Path      string    `json:"path"`
	DeviceID  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
}
