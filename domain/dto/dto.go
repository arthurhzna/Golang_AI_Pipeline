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

type PredictionMQTTPayload struct {
	DeviceID           string  `json:"device_id"`
	Timestamp          string  `json:"timestamp"`
	OutputText         string  `json:"output_text"`
	OutputPath         string  `json:"output_path"`
	PredictedPlatColor string  `json:"predicted_plat_color"`
	PredictedPlatType  string  `json:"predicted_plat_type"`
	TimeTakenPredict   float64 `json:"time_taken_predict"`
}
