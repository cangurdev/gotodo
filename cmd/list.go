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
		fmt.Print("You have the following tasks:\n\n")
		color.Cyan("Id\tStatus\tDue\tFolder\tTask")
		for i, task := range tasks {
			task := task.Value
			line := fmt.Sprintf("%d\t%v\t%s\t%s\t%s\n", i+1, task.IsDone, task.Due, task.Parent, task.Text)

			if !task.IsDone {
				if task.Tag == "important" {
					color.Red(line)
				} else {
					fmt.Printf(line)
				}
			} else {
				color.Green("%d.\t%v\t%s\t%s\t%s\n", i+1, task.IsDone, task.Due, task.Parent, task.Text)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
