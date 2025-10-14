package cmd

import (
	"os"
	"task_queue/config"

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

	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
