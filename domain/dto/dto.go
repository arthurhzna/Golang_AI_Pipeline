package dto

import (
	"mime/multipart"
)

type QueueRequest struct {
	Images    *multipart.FileHeader `form:"image" binding:"required"`
	DeviceId  string                `form:"device_id" binding:"required"`
	Timestamp string                `form:"timestamp" binding:"required"`
}

type QueueResponse struct {
	DeviceId  string `json:"device_id"`
	Timestamp string `json:"timestamp"`
}
