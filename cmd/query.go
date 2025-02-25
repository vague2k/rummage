package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
)

// TODO: add docs for new functionality / tests
func newQueryCmd(db database.RummageDbInterface) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Query the database to find an entry by highest score, or using an exact match",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				commands.Query(cmd, arg, db)
			}
		},
	}

	queryCmd.Flags().IntP("quantity", "p", 10, "")

	return queryCmd
}
