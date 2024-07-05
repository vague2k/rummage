package populate

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vague2k/rummage/internal"
	"github.com/vague2k/rummage/utils"
)

var logger = internal.NewLogger(nil, os.Stdout)

// Walks $GOPATH/pkg/mod/github.com and parse out package paths valid as args in a "go get" command.
//
// Returns a []strings where each elem looks like "github.com/gorilla/mux" for example.
//
// TODO: write test for this func
func WalkAndParsePackages(dir string) []string {
	var (
		pkgs []string
		seen = make(map[string]bool)
	)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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
	if err != nil {
		logger.Fatal("Could not walk dirs: \n", err)
	}

	return pkgs
}
