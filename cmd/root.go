package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/internal"
)

var logger = internal.NewLogger(nil, os.Stdout)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "rummage",
	Version: "2.0.0-beta",
	Short:   "A zoxide inspired alternative to go get",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		populateCmd,
		removeCmd,
		getCmd,
		addCmd,
		queryCmd,
	)
}
