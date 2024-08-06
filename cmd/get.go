package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
)

func newGetCmd(db database.RummageDbInterface) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a go package from the database using a substring, or get a package how you normally would",
		Run: func(cmd *cobra.Command, args []string) {
			commands.Get(cmd, args, db)
		},
	}

	getCmd.Flags().BoolP("update", "u", false, "same as '-u', see 'go help get'")
	getCmd.Flags().BoolP("dependencies", "t", false, "same as '-t', see 'go help get'")
	getCmd.Flags().BoolP("debug", "x", false, "same as '-x', see 'go help get'")

	return getCmd
}
