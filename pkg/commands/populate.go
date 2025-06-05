package commands

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

// package version (e.g. @v1.0.0) will be stripped, but packages such as "github.com/user/example/v2" are valid
func cut(s string) string {
	_, after, _ := strings.Cut(s, "mod/")
	pkg, _, _ := strings.Cut(after, "@")
	return pkg
}

func capitalize(author string) string {
	var result []rune
	i := 0
	for i < len(author) {
		if author[i] == '!' && i+1 < len(author) {
			result = append(result, rune(author[i+1]-32)) // convert to uppercase
			i += 2
		} else {
			result = append(result, rune(author[i]))
			i++
		}
	}
	return string(result)
}

// Walks a dir and extracts packages valid as args in a "go get" command.
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

		curr := cut(path)
		if strings.Count(curr, "/") < 2 {
			return nil
		}

		// keep track of items added to the slice, instead of walked
		if seen[curr] {
			return nil
		}
		seen[curr] = true

		normalized := capitalize(curr)

		pkgs = append(pkgs, normalized)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could walk directory %s to extract packages", dir)
	}

	return pkgs, nil
}

// The "populate" command populates the database with third party packages already known by go.
// By the default, the dir "populate" walks through will be "$GOPATH/pkg/mod/github.com"
func Populate(cmd *cobra.Command, args []string, db *database.Queries, ctx context.Context) {
	flagDir := cmd.Flag("dir").Value.String()

	pkgs, err := extractPackages(flagDir)
	if err != nil {
		cmd.PrintErr(err)
	}

	amtAdded := 0
	for _, entry := range pkgs {
		_, err := db.AddItem(ctx, database.AddItemParams{
			Entry:        entry,
			Score:        1.0,
			Lastaccessed: time.Now().Unix(),
		})
		if err != nil {
			cmd.PrintErr(err)
		}
		amtAdded++
	}

	if amtAdded == 0 {
		cmd.PrintErrf("no new packages were found to populate the database, added %d packages\n", amtAdded)
	} else {
		cmd.Printf("added %d packages\n", amtAdded)
	}
}
