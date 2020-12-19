package cmd

import "github.com/spf13/cobra"

//RootCmd function
var RootCmd = &cobra.Command{
	Use:   "gotodo",
	Short: "GoToDo is a CLI task manager",
}
