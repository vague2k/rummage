package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

func NewRootCmd(db database.RummageDbInterface) *cobra.Command {
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

	rootCmd.AddCommand(newPopulateCmd(db))
	rootCmd.AddCommand(newQueryCmd(db))
	rootCmd.AddCommand(newAddCmd(db))
	rootCmd.AddCommand(newRemoveCmd(db))

	return rootCmd
}
