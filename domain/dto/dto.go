package dto

import (
	"mime/multipart"
	"time"
)

type UploadRequest struct {
	Images    *multipart.FileHeader `form:"file" binding:"required"`
	Timestamp time.Time             `form:"timestamp" binding:"required"`
	DeviceId  string                `form:"device_id" binding:"required"`
}

type UploadResponse struct {
	DeviceId  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
}
