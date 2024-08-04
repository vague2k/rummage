package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
)

func newRemoveCmd(db database.RummageDbInterface) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove items from the database or be prompted to confirm to remove all items",
		Run: func(cmd *cobra.Command, args []string) {
			commands.Remove(cmd, args, db)
		},
	}

	removeCmd.Flags().BoolP("delete-all", "D", false, "Remove all items from the database, you will be prompted to confirm this choice")

	return removeCmd
}
