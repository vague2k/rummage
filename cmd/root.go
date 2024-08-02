package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "rummage [command]",
		Version: "3.0.0",
		Short:   "A zoxide inspired alternative to go get",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				cmd.PrintErr(err)
			}
		},
	}

	return rootCmd
}
