package dto

import (
	"mime/multipart"
	"time"
)

type QueueRequest struct {
	Images    *multipart.FileHeader `form:"image" binding:"required"`
	DeviceId  string                `form:"device_id" binding:"required"`
	Timestamp time.Time             `form:"timestamp" binding:"required"`
}

type QueueResponse struct {
	DeviceId  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
}
