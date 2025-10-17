package models

type QueueDataRedis struct {
	FileName  string `json:"file_name"`
	Path      string `json:"file_path"`
	DeviceID  string `json:"device_id"`
	Timestamp string `json:"timestamp"`
}

type QueuePredictionRedis struct {
	DeviceID           string  `json:"device_id"`
	Timestamp          string  `json:"timestamp"`
	FileName           string  `json:"file_name"`
	ImageOutputPath    string  `json:"image_output_path"`
	OutputText         string  `json:"output_text"`
	PredictedPlatColor string  `json:"predicted_plat_color"`
	PredictedPlatType  string  `json:"predicted_plat_type"`
	TimeTakenPredict   float64 `json:"time_taken_predict"`
}
