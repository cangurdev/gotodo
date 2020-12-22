package cmd

import (
	"fmt"
	"strings"

	"github.com/cvngur/gotodo/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:              "add",
	TraverseChildren: true,
	Short:            "Adds a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		importanTag, _ := cmd.Flags().GetBool("important")
		due, _ := cmd.Flags().GetString("due")
		parent, _ := cmd.Flags().GetString("folder")
		task := strings.Join(args, " ")

		err := db.CreateTask(task, due, parent, importanTag)
		if err != nil {
			fmt.Println("Something went wrong:", err)
		}
		fmt.Printf("Added '%s' to your task list\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("important", "i", false, "Important task")
	addCmd.Flags().StringP("due", "d", "-", "Due time of the task")
	addCmd.Flags().StringP("folder", "f", "-", "Parent folder of the task")
}
