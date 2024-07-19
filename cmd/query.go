package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the database to find a record's details",
	Run: func(cmd *cobra.Command, args []string) {
		flags := utils.RegisterBoolFlags(cmd, "exact")

		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		switch true {
		case flags["exact"]:
			for _, arg := range args {
				commands.FindExactMatch(db, arg)
			}
		default:
			for _, arg := range args {
				commands.FindHighestScore(db, arg)
			}
		}
	},
}

func init() {
	queryCmd.Flags().BoolP("exact", "e", false, "Query using an exact substring match instead of using highest score.")
}

