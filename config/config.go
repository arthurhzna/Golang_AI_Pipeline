package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Config AppConfig

type AppConfig struct {
	AppPort string

	RedisAddr     string
	RedisPassword string
	KeyRedisSend  string
	KeyRedisGet   string

	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSBucket          string
	AWSPathBucket      string

	MQTTBroker   string
	MQTTPort     string
	MQTTClientID string
	MQTTUsername string
	MQTTPassword string
	MQTTTopic    string

	NumWorkers int

	BaseDirSend string
	BaseDirGet  string

	APIKey string
}

func Init() {
	_ = godotenv.Load()

	numWorkers, err := strconv.Atoi(os.Getenv("WORKER"))
	if err != nil {
		numWorkers = 1
	}

	Config = AppConfig{
		AppPort: os.Getenv("APP_PORT"),

		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		KeyRedisSend:  os.Getenv("KEY_REDIS_SEND"),
		KeyRedisGet:   os.Getenv("KEY_REDIS_GET"),

		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_DEFAULT_REGION"),
		AWSBucket:          os.Getenv("AWS_BUCKET"),
		AWSPathBucket:      os.Getenv("AWS_PATH_BUCKET"),

		MQTTBroker:   os.Getenv("MQTT_BROKER"),
		MQTTPort:     os.Getenv("MQTT_PORT"),
		MQTTClientID: os.Getenv("CLIENT_MQTT_ID"),
		MQTTUsername: os.Getenv("MQTT_USERNAME"),
		MQTTPassword: os.Getenv("MQTT_PASSWORD"),
		MQTTTopic:    os.Getenv("MQTT_TOPIC"),

		NumWorkers: numWorkers,

		BaseDirSend: os.Getenv("BASE_DIR_SEND"),
		BaseDirGet:  os.Getenv("BASE_DIR_GET"),

		APIKey: os.Getenv("API_KEY"),
	}

}
