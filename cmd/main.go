package cmd

import (
	"context"
	"task_queue/config"
	"task_queue/controllers"
	"task_queue/repositories"
	"task_queue/services"

	"fmt"
	"net/http"
	"task_queue/common/aws"
	"task_queue/common/mqtt"
	"task_queue/common/response"
	"task_queue/constants"
	"task_queue/middlewares"
	"task_queue/routes"
	"task_queue/workers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "task-queue",
	Short: "task queue",
	Long:  "task queue",
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		fmt.Println(config.Config)
		redisClient, err := config.CreateClient(context.Background(), 0, config.Config.RedisAddr, config.Config.RedisPassword)
		if err != nil {
			panic(err)
		}

		awsS3 := aws.NewAWS_S3(
			config.Config.AWSAccessKeyID,
			config.Config.AWSSecretAccessKey,
			config.Config.AWSRegion,
			config.Config.AWSBucket,
		)

		err = awsS3.CreateClient(context.Background())
		if err != nil {
			panic(err)
		}

		mqtt := mqtt.NewMQTT(
			config.Config.MQTTBroker,
			config.Config.MQTTPort,
			config.Config.MQTTClientID,
			config.Config.MQTTUsername,
			config.Config.MQTTPassword,
		)

		err = mqtt.Connect(context.Background())
		if err != nil {
			panic(err)
		}

		queueRepository := repositories.NewQueueRepository(
			redisClient,
			config.Config.KeyRedisSend,
			config.Config.KeyRedisGet,
		)

		queueService := services.NewQueueService(
			awsS3,
			mqtt,
			config.Config.AWSPathBucket,
			config.Config.MQTTTopic,
			queueRepository,
			config.Config.BaseDirSend,
			config.Config.BaseDirGet,
		)

		queueController := controllers.NewQueueController(
			queueService,
		)

		worker := workers.NewWorker(
			queueService,
		)

		err = worker.Run(context.Background(), config.Config.NumWorkers)
		if err != nil {
			panic(err)
		}

		router := gin.Default()
		router.Use(middlewares.HandlePanic())

		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to Task Queue",
			})
		})
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-api-key")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		group := router.Group("/task-queue")
		route := routes.NewRoute(queueController, group)
		route.Serve()
		router.Run(config.Config.AppPort)
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
