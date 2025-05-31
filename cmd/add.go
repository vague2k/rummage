package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
)

func newAddCmd(db *database.Queries, ctx context.Context) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add an item that resembles a go package manually to the database",
		Run: func(cmd *cobra.Command, args []string) {
			commands.Add(cmd, args, db, ctx)
		},
	}

	return addCmd
}
