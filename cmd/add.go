package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	addCmd.Flags().BoolP("multiple", "m", false, "Enables the adding of multiple items.")
	addCmd.Flags().BoolP("stdout", "o", false, "Set the added item's details to stdout.")
	addCmd.Flags().BoolP("print", "p", false, "Prints the added item's details.")
}
