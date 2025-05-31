package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

// Add allows a user to manually add items to the database
func Add(cmd *cobra.Command, args []string, db *database.Queries, ctx context.Context) {
	if len(args) == 1 {
		if err := utils.ResemblesGoPackage(args[0]); err != nil {
			cmd.PrintErrf("%s\n", err)
			return
		}
		_, err := db.AddItem(ctx, database.AddItemParams{
			Entry: args[0],
		})
		if err != nil {
			cmd.PrintErrf("%s\n", err)
			return
		}
		cmd.Printf("added 1 package\n")
		return
	}

	amtAdded := 0
	for _, entry := range args {
		if err := utils.ResemblesGoPackage(entry); err != nil {
			issue := fmt.Sprintf("issue occured when adding item %s to the db\n", entry)
			added := fmt.Sprintf("added %d packages\n", amtAdded)
			cmd.PrintErrf("%s%s%s",
				"stopping the adding of items early...\n",
				issue,
				added)
			return
		}
		_, err := db.AddItem(ctx, database.AddItemParams{
			Entry: entry,
		})
		if err != nil {
			cmd.PrintErr("stopping the adding of items early...\n")
			cmd.PrintErrf("%s\n", err)
		}
		amtAdded++
	}

	cmd.Printf("added %d packages\n", amtAdded)
}
