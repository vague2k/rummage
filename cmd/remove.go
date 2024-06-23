package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/pkg/ui"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a single item or multiple items from the database. Must be a full package path",
	Run: func(cmd *cobra.Command, args []string) {
		flagDeleteAll, err := cmd.Flags().GetBool("delete-all")
		if err != nil {
			logger.Fatal("Could not get '--delete-all' flag value")
		}
		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		if flagDeleteAll {
			choice := ui.YesNoPrompt("Are you sure you want to delete all items from the database?", false)
			if !choice {
				return
			}
			err := db.DeleteAllItems()
			if err != nil {
				logger.Fatal(err)
			}
			logger.Warn("Deleted ALL from the database...")
		}

		for _, arg := range args {
			if _, err := db.DeleteItem(arg); err != nil {
				logger.Fatal(err)
			}
		}
	},
}

func init() {
	removeCmd.Flags().BoolP("delete-all", "D", false, "Delete all items from the database, confirm with a prompt")
}
