package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/internal"
	"github.com/vague2k/rummage/pkg"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

var logger = internal.NewLogger(nil, os.Stdout)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "rummage",
	Version: LatestVersion(),
	Short:   "A zoxide inspired alternative to go get",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		for _, arg := range args {
			hasSlash, _ := utils.ParseForwardSlash(arg)
			// can safely assume if a "/" is parsed from the arg, it's more than likely an absolute pkg path
			if hasSlash {
				pkg.GoGetAddedItem(db, arg)
				return
			}
			pkg.GoGetHighestScore(db, arg)
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
	rootCmd.AddCommand(quickstartCmd)
}
