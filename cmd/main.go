package cmd

import (
	"context"
	"os"
	"strconv"
	"task_queue/config"
	"task_queue/controllers"
	"task_queue/repositories"
	"task_queue/services"

	"fmt"
	"net/http"
	"task_queue/common/aws"
	"task_queue/common/response"
	"task_queue/constants"
	"task_queue/middlewares"
	"task_queue/routes"
	"task_queue/workers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "task-queue",
	Short: "task queue",
	Long:  "task queue",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()
		redisClient, err := config.CreateClient(context.Background(), 0, os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"))
		if err != nil {
			panic(err)
		}

		queueRepository := repositories.NewQueueRepository(redisClient, os.Getenv("KEY_REDIS_SEND"), os.Getenv("KEY_REDIS_GET"))
		queueService := services.NewQueueService(queueRepository, os.Getenv("BASE_DIR_SEND"), os.Getenv("BASE_DIR_GET"))
		queueController := controllers.NewQueueController(queueService)

		awsS3 := aws.NewAWS_S3(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_REGION"), os.Getenv("AWS_BUCKET"))
		if err != nil {
			panic(err)
		}

		worker := workers.NewWorker(queueRepository, awsS3)
		numWorkers, err := strconv.Atoi(os.Getenv("WORKER"))
		if err != nil {
			numWorkers = 1
		}
		err = worker.Run(context.Background(), numWorkers)
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
		router.Run(os.Getenv("APP_PORT"))
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
