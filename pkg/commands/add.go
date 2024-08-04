package commands

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

// Add allows a user to manually add items to the database
func Add(cmd *cobra.Command, args []string, db database.RummageDbInterface) {
	if len(args) == 1 {
		_, err := db.AddItem(args[0])
		if err != nil {
			cmd.PrintErrf("%s\n", err)
			return
		}
		cmd.Printf("added 1 package\n")
		return
	}

	_, amtAdded, err := db.AddMultiItems(args...)
	if err != nil {
		cmd.PrintErr("stopping the adding of items early...\n")
		cmd.PrintErrf("%s\n", err)
	}

	cmd.Printf("added %d packages\n", amtAdded)
}
