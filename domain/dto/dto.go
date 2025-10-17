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
	Timestamp_In       string  `json:"timestamp_In"`
	Timestamp_Out      string  `json:"timestamp_Out"`
	FileName           string  `json:"file_name"`
	ImageOutputPath    string  `json:"image_output_path"`
	OutputText         string  `json:"output_text"`
	PredictedPlatColor string  `json:"predicted_plat_color"`
	PredictedPlatType  string  `json:"predicted_plat_type"`
	TimeTakenPredict   float64 `json:"time_taken_predict"`
}
