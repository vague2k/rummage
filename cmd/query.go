package cmd

import (
	"github.com/spf13/cobra"
	s "github.com/vague2k/rummage/cmd/services/query"
	"github.com/vague2k/rummage/pkg/database"
	cmdUtils "github.com/vague2k/rummage/utils/cmd"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the database to find a record's details",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmdUtils.RegisterBoolFlags(cmd, "exact")

		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		switch true {
		case flags["exact"]:
			for _, arg := range args {
				s.FindExactMatch(db, arg)
			}
		default:
			for _, arg := range args {
				s.FindHighestScore(db, arg)
			}
		}
	},
}

func init() {
	queryCmd.Flags().BoolP("exact", "e", false, "Query using an exact substring match instead of using highest score.")
}