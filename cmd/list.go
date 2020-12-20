package cmd

import (
	"fmt"
	"os"

	"github.com/cvngur/gotodo/db"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists tasks in your task list",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong")
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks")
			return
		}
		fmt.Println("You have the following tasks:\n")
		for i, task := range tasks {
			isDone := task.Value.IsDone
			line := fmt.Sprintf("%d. %s.\n", i+1, task.Value.Text)
			tag := task.Value.Tag
			if !isDone {
				if tag == "important" {
					color.Red(line)
				} else if tag == "today" {
					color.Yellow(line)
				} else {
					fmt.Printf(line)
				}
			}
		}
		fmt.Println("\nCompleted tasks:\n")
		for j, task := range tasks {
			isDone := task.Value.IsDone
			if isDone {
				color.Green("%d. %s.\n", j+1, task.Value.Text)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
