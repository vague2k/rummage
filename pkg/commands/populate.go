package commands

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

// Cuts a path string into a proper, valid go package that is able to be got.
//
// package version (e.g. @v1.0.0) will be stripped, but packages such as "github.com/user/example/v2" are valid
func cut(path string) string {
	_, after, _ := strings.Cut(path, "mod/")
	pkg, _, _ := strings.Cut(after, "@")
	return pkg
}

// Walks a dir (should be $GOPATH/pkg/mod/github.com) and extracts packages valid as args in a "go get" command.
func extractPackages(dir string) ([]string, error) {
	var pkgs []string
	seen := make(map[string]bool)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		pkg := cut(path)
		// if the package path does not have at least 2 matches of "/", it's not a valid go package, and can be skipped
		if strings.Count(pkg, "/") < 2 {
			return nil
		}
		// keep track of items added to the slice, instead of walked
		if seen[pkg] {
			return nil
		}
		seen[pkg] = true
		pkgs = append(pkgs, pkg)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

// Populate will fill the database with packages already installed from you go path's "mod/github.com" directory
func Populate(db *database.RummageDB) {
	GOPATH := utils.UserGoPath()
	dir := filepath.Join(GOPATH, "pkg", "mod", "github.com")
	pkgs, err := extractPackages(dir)
	if err != nil {
		log.Fatal(err)
	}

	items, err := db.AddMultiItems(pkgs...)
	if err != nil {
		log.Fatal("Failure adding items: \n", err)
	}

	log.Print(fmt.Sprintf("Added %d packages to the database.\n", len(items)))
}
