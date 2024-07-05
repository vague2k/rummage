package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item to the database",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		if _, err := db.AddMultiItems(args...); err != nil {
			logger.Fatal(err)
		}
		return
	},
}
