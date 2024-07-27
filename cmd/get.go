package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a go package from the database, and increment it's recency score",
	Run: func(cmd *cobra.Command, args []string) {
		flags := utils.RegisterBoolFlags(cmd, "update", "test-deps", "debug")
		db, err := database.Init("")
		if err != nil {
			log.Fatal(err)
		}

		var goGetFlags []string
		if flags["update"] {
			goGetFlags = append(goGetFlags, "-u")
		} else if flags["test-deps"] {
			goGetFlags = append(goGetFlags, "-t")
		} else if flags["debug"] {
			goGetFlags = append(goGetFlags, "-x")
		} else {
			goGetFlags = nil
		}

		for _, arg := range args {
			if err := commands.Get(db, arg, goGetFlags...); err != nil {
				log.Err(err)
				return
			}
		}
	},
}

func init() {
	getCmd.Flags().BoolP("update", "u", false, "Same thing as '-u'. See 'go help get'")
	getCmd.Flags().BoolP("test-deps", "t", false, "Same thing as '-t'. See 'go help get'")
	getCmd.Flags().BoolP("debug", "x", false, "Same thing as '-x'. See 'go help get'")
}
