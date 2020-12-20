package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cvngur/gotodo/db"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task in your task list",
	Run: func(cmd *cobra.Command, args []string) {

		key, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Failed")
			os.Exit(1)
		}
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong", err)
			return
		}
		if key <= 0 || key > len(tasks) {
			fmt.Println("Invalid task number: ", key)
			os.Exit(1)
		}
		task := tasks[key-1]
		newTask := strings.Join(args[1:], " ")

		err = db.UpdateTask(task.Key, newTask)
		if err != nil {
			fmt.Println("Something went wrong:", err)
		}
		fmt.Println("Updated task no:", key)
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
