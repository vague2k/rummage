package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Add already all installed packages to the database",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Grabbing manually installed packages to add to the database...")
		db, err := database.Init("")
		if err != nil {
			logger.Fatal(err)
		}
		defer db.DB.Close()

		GOPATH := utils.UserGoPath()
		dir := filepath.Join(GOPATH, "pkg/mod/github.com")

		var (
			pkgs []string
			seen = make(map[string]bool)
		)
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				logger.Err(err)
				return nil
			}
			if !d.IsDir() {
				return nil
			}
			// isolate a proper package path, (i.e "github.com/gorilla/mux")
			// "github.com/gorilla/mux" and "github.com/gorilla/mux/v2" are valid
			// adds latest versions since version # is stripped
			pkg := strings.Split(path, "mod/")
			pkgNoVerNum := strings.Split(pkg[1], "@")
			curr := pkgNoVerNum[0]
			// if the package path does not have at least 2 matches of "/", it's not a valid go package, and can be skipped
			_, slashes := utils.ParseForwardSlash(curr)
			if slashes < 2 {
				return nil
			}

			if seen[curr] {
				return nil
			}
			seen[curr] = true

			pkgs = append(pkgs, pkgNoVerNum[0])
			return nil
		})

		items, err := db.AddMultiItems(pkgs...)
		if err != nil {
			logger.Fatal("Failure adding items: \n", err)
		}

		msg := fmt.Sprintf("All done. Added %d go packages.", len(items))
		logger.Info(msg)
	},
}
