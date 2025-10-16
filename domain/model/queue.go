package models

type QueueDataRedis struct {
	Path      string `json:"path"`
	DeviceID  string `json:"device_id"`
	Timestamp string `json:"timestamp"`
}

type QueuePredictionRedis struct {
	DeviceID           string  `json:"device_id"`
	Timestamp          string  `json:"timestamp"`
	OutputText         string  `json:"output_text"`
	OutputPath         string  `json:"output_path"`
	PredictedPlatColor string  `json:"predicted_plat_color"`
	PredictedPlatType  string  `json:"predicted_plat_type"`
	TimeTakenPredict   float64 `json:"time_taken_predict"`
}
