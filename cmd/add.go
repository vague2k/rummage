package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item to the database",
	Run: func(cmd *cobra.Command, args []string) {
		flagMultiple, err := cmd.Flags().GetBool("multiple")
		if err != nil {
			logger.Fatal("Could not get '--multiple' flag value")
		}

		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		if flagMultiple {
			if _, err := db.AddMultiItems(args...); err != nil {
				logger.Fatal(err)
			}
			return
		}

		if _, err := db.AddItem(args[0]); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	addCmd.Flags().BoolP("multiple", "m", false, "Add multiple items to the database")
}
