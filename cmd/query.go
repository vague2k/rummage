package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the database to find a record's details",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}
		for _, arg := range args {
			arg = strings.ToLower(arg)
			found, exists := db.EntryWithHighestScore(arg)
			if !exists {
				logger.Warn("No entry was found that could match your query.")
				return
			}
			info := fmt.Sprintf("Entry: %s\nScore: %f\nLastAccessed: %d\n",
				found.Entry,
				found.Score,
				found.LastAccessed,
			)
			logger.Info("found top match for ", arg, "\n", info)
		}
	},
}
