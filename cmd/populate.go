package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
)

var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "Add already all installed packages to the database",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Init("")
		if err != nil {
			log.Fatal(err)
		}
		defer db.DB.Close()

		commands.Populate(db)
	},
}
