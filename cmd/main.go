package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "task-queue",
	Short: "task queue",
	Long:  "task queue",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
