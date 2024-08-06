package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

func newPopulateCmd(db database.RummageDbInterface) *cobra.Command {
	populateCmd := &cobra.Command{
		Use:   "populate",
		Short: "Populate the database with third party packages already known by go",
		Run: func(cmd *cobra.Command, args []string) {
			commands.Populate(cmd, args, db)
		},
	}

	populateCmd.Flags().StringP("dir", "d", filepath.Join(utils.UserGoPath(), "pkg", "mod", "github.com"), "The directory you want to search for valid go packages")

	return populateCmd
}
