package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
	cmdUtils "github.com/vague2k/rummage/utils/cmd"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a single item or multiple items from the database. Must be a full package path",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmdUtils.RegisterBoolFlags(cmd, "delete-all")

		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		switch true {
		case flags["delete-all"]:
			commands.PromptDeleteAll(db)
		default:
			for _, arg := range args {
				if _, err := db.DeleteItem(arg); err != nil {
					logger.Fatal(err)
				}
			}
		}
	},
}

func init() {
	removeCmd.Flags().BoolP("delete-all", "D", false, "Delete all items from the database, confirm with a prompt")
}
