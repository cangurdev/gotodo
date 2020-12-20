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
		var doneTasks []string
		if err != nil {
			fmt.Println("Something went wrong")
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks")
			return
		}
		fmt.Print("You have the following tasks:\n\n")
		for i, task := range tasks {
			task := task.Value
			line := fmt.Sprintf("%d. %s.\n", i+1, task.Text)

			if !task.IsDone {
				if task.Tag == "important" {
					color.Red(line)
				} else if task.Tag == "today" {
					color.Yellow(line)
				} else {
					fmt.Printf(line)
				}
			} else {
				doneTasks = append(doneTasks, task.Text)
			}
		}
		fmt.Print("\nCompleted tasks:\n\n")
		for j, task := range doneTasks {
			color.Green("%d. %s.\n", j+1, task)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
