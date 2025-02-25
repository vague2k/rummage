package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/config"
	"github.com/vague2k/rummage/pkg/database"
)

func NewRootCmd(db database.RummageDbInterface) *cobra.Command {
	conf := config.SetVersions()
	rootCmd := &cobra.Command{
		Use:     "rummage [command]",
		Version: conf.Rummage.Version,
		Short:   "A smart wrapper around 'go get'",
		Run: func(cmd *cobra.Command, args []string) {
			dbVersionFlag, err := cmd.Flags().GetBool("apiver")
			if err != nil {
				panic(err)
			}

			if dbVersionFlag {
				cmd.Printf("rummage database api version %s\n", db.Version())
				return
			}

			if err := cmd.Help(); err != nil {
				cmd.PrintErr(err)
			}
		},
	}

	rootCmd.Flags().BoolP("apiver", "a", false, "outputs the current version of the rummage database api")

	rootCmd.AddCommand(newPopulateCmd(db))
	rootCmd.AddCommand(newQueryCmd(db))
	rootCmd.AddCommand(newAddCmd(db))
	rootCmd.AddCommand(newRemoveCmd(db))
	rootCmd.AddCommand(newGetCmd(db))

	return rootCmd
}
