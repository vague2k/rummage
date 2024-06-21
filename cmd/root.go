package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/internal"
	"github.com/vague2k/rummage/pkg"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/pkg/utils"
)

var logger = internal.NewLogger(nil, os.Stdout)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rummage",
	Short: "A zoxide inspired alternative to go get",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		for _, arg := range args {
			hasSlash := utils.ParseForwardSlash(arg)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rummage.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rootCmd.AddCommand(addCmd)
}
