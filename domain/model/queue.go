package models

type QueueDataRedis struct {
	Path      string `json:"path"`
	DeviceID  string `json:"device_id"`
	Timestamp string `json:"timestamp"`
}
