package cmd

import (
	"os"
	"task_queue/config"
	"task_queue/controllers"
	"task_queue/repositories"
	"task_queue/services"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "task-queue",
	Short: "task queue",
	Long:  "task queue",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()
		redisClient, err := config.CreateClient(0, os.Getenv("DB_ADDR"), os.Getenv("DB_PASS"))
		if err != nil {
			panic(err)
		}

		queueRepository := repositories.NewQueueRepository(redisClient)
		queueService := services.NewQueueService(queueRepository)
		queueController := controllers.NewQueueController(queueService)

		// router := gin.Default()
		// router.Use(middlewares.HandlePanic())

	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
