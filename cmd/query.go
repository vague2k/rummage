package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
)

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

	queryCmd.Flags().BoolP("exact", "e", false, "Query the database to find an entry using an exact search, instead of by highest score")
	queryCmd.Flags().BoolP("score", "s", false, "Display the query's score")
	queryCmd.Flags().BoolP("last-accessed", "l", false, "Display the last accessed timestamp of the query in Unix seconds")

	return queryCmd
}
