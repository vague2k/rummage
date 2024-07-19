package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "Add already all installed packages to the database",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Grabbing manually installed packages to add to the database...")
		db, err := database.Init("")
		if err != nil {
			log.Fatal(err)
		}
		defer db.DB.Close()

		GOPATH := utils.UserGoPath()
		dir := filepath.Join(GOPATH, "pkg", "mod", "github.com")

		pkgs := commands.WalkAndParsePackages(dir)

		items, err := db.AddMultiItems(pkgs...)
		if err != nil {
			log.Fatal("Failure adding items: \n", err)
		}

		msg := fmt.Sprintf("All done. Added %d go packages.", len(items))
		log.Info(msg)
	},
}
