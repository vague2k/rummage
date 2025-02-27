package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands/utils"
	"github.com/vague2k/rummage/pkg/config"
	"github.com/vague2k/rummage/pkg/database"
)

func NewRootCmd(db database.RummageDbInterface) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "rummage [command]",
		Short: "A smart wrapper around 'go get'",
		Run: func(cmd *cobra.Command, args []string) {
			flags, err := utils.GetBoolFlags(cmd, "version", "apiver")
			if err != nil {
				panic(err)
			}

			// the "Version" field in the cobra.Command struct can't be used because
			// that would cause me to pass in "config" as a param affecting tests.
			//
			// right now it's not worth refactoring all that code when I could just do this, and just mark both flags as mutually exclusive
			if flags["apiver"] {
				cmd.Printf("rummage database api version %s\n", db.Version())
				return
			}

			ver := config.SetVersions().Rummage.Version

			if flags["version"] {
				cmd.Printf("rummage version %s\n", ver)
				return
			}

			if err := cmd.Help(); err != nil {
				cmd.PrintErr(err)
			}
		},
	}

	rootCmd.Flags().BoolP("apiver", "a", false, "outputs the current version of the rummage database api")
	rootCmd.Flags().BoolP("version", "v", false, "outputs the current rummage version")
	rootCmd.MarkFlagsMutuallyExclusive("version", "apiver")

	rootCmd.AddCommand(newPopulateCmd(db))
	rootCmd.AddCommand(newQueryCmd(db))
	rootCmd.AddCommand(newAddCmd(db))
	rootCmd.AddCommand(newRemoveCmd(db))
	rootCmd.AddCommand(newGetCmd(db))

	return rootCmd
}
