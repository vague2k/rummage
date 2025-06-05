package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

func NewRootCmd(db *database.Queries) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "rummage [command]",
		Version: "v3.3.0",
		Short:   "A smart wrapper around 'go get'",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				cmd.PrintErr(err)
			}
		},
	}

	ctx := context.Background()
	rootCmd.AddCommand(newPopulateCmd(db, ctx))
	rootCmd.AddCommand(newQueryCmd(db, ctx))
	rootCmd.AddCommand(newAddCmd(db, ctx))
	rootCmd.AddCommand(newRemoveCmd(db, ctx))
	rootCmd.AddCommand(newGetCmd(db, ctx))

	return rootCmd
}
