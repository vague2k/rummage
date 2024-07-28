package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item to the database",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Init("")
		if err != nil {
			log.Fatal(err)
		}

		err = commands.Add(db, args...)
		if err != nil {
			log.Err(err)
		}
	},
}
