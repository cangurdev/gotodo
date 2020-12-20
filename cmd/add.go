package cmd

import (
	"fmt"
	"strings"

	"github.com/cvngur/gotodo/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:              "add [OPTIONS]",
	TraverseChildren: true,
	Short:            "Adds a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		importanTag, _ := cmd.Flags().GetBool("important")
		todayTag, _ := cmd.Flags().GetBool("today")
		task := strings.Join(args, " ")
		tag := ""
		if importanTag {
			tag = "important"
		} else if todayTag {
			tag = "today"
		}
		err := db.CreateTask(task, tag)
		if err != nil {
			fmt.Println("Something went wrong:", err)
		}
		fmt.Printf("Added '%s' to your task list\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("important", "i", false, "Important task")
	addCmd.Flags().BoolP("today", "t", false, "Task should done today")
}
