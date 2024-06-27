package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/cmd/services"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a go package from the database, and increment it's recency score",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}

		for _, arg := range args {
			hasSlash, _ := utils.ParseForwardSlash(arg)
			// can safely assume if a "/" is parsed from the arg, it's more than likely an absolute pkg path
			if hasSlash {
				services.GoGetAddedItem(db, arg)
				return
			}
			services.GoGetHighestScore(db, arg)
		}
	},
}
