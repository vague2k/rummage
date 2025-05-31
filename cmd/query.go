package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
)

func newQueryCmd(db *database.Queries, ctx context.Context) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Query the database to find an entry by highest score, or using an exact match",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				commands.Query(cmd, arg, db, ctx)
			}
		},
	}

	queryCmd.Flags().IntP("quantity", "q", 10, "The amount of entry matches to display in the output")

	return queryCmd
}
