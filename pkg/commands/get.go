package commands

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands/utils"
	"github.com/vague2k/rummage/pkg/database"
)

// attempts to call "go get" against an arg and also takes into account any flags "go get" may accept
func goGet(cobra *cobra.Command, pkg string, flags ...string) error {
	cmd := exec.Command("go", "get")
	if flags == nil && pkg != "" {
		cmd.Args = append(cmd.Args, pkg)
	}
	if pkg == "" {
		cmd.Args = append(cmd.Args, flags...)
	} else {
		cmd.Args = append(cmd.Args, flags...)
		cmd.Args = append(cmd.Args, pkg)
	}

	b, err := cmd.CombinedOutput()
	output := string(b)
	if err != nil {
		return fmt.Errorf("%s", output)
	}

	cobra.Print(output)
	return nil
}

// "go get"s a package and adds it to the database, and update it's score
//
// it's assumed that if this function is called, the item does not yet exist in the database
func getAddedItem(cmd *cobra.Command, db *database.Queries, ctx context.Context, pkg string, flags ...string) {
	item, err := db.AddItem(ctx, database.AddItemParams{
		Entry: pkg,
	})
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}

	err = goGet(cmd, pkg, flags...)
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}

	err = db.UpdateItem(ctx, database.UpdateItemParams{
		Entry:        item.Entry,
		Score:        item.RecalculateScore(),
		Lastaccessed: time.Now().Unix(),
	})
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}
}

// "go get"s a package based on a db search and score, and at the end, update it's score
//
// it's assumed that if this function is called, the item exists in the database
func getHighestScore(cmd *cobra.Command, db *database.Queries, ctx context.Context, pkgSubstr string, flags ...string) {
	item, err := db.EntryWithHighestScore(ctx, "%"+pkgSubstr+"%")
	if err != nil && err == sql.ErrNoRows {
		cmd.PrintErrf("no match found with the given arguement %s\n", pkgSubstr)
		return
	}
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}

	err = goGet(cmd, item.Entry, flags...)
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}

	err = db.UpdateItem(ctx, database.UpdateItemParams{
		Entry:        item.Entry,
		Score:        item.RecalculateScore(),
		Lastaccessed: time.Now().Unix(),
	})
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}
}

// helper function to get flags that could be used in a "go get" call
func getFlags(flagMap map[string]bool) []string {
	var flags []string
	if flagMap["update"] {
		flags = append(flags, "-u")
	} else if flagMap["dependencies"] {
		flags = append(flags, "-t")
	} else if flagMap["debug"] {
		flags = append(flags, "-x")
	} else {
		flags = nil
	}
	return flags
}

func handleFlagsNoPkg(cmd *cobra.Command, flags ...string) {
	err := goGet(cmd, "", flags...)
	if err != nil {
		cmd.PrintErrf("%s\n", err)
	}
}

// the "get" command attempts to get a package based on a couple criteria.
//
// If the arguement (not flag) passed to "get" has 2 slashes or more, it's assumed the item does not yet exist
// in the database and will add it to the database while also calling "go get" upon that arguement
//
// If the arguement (not flag) passed looks more like a substring (e.g rummage get mux) then it's assumed the item
// exists in the database and a search based on highest score will be performed on that arguement.
//
// In both cases if a match for an arguement can't be found, the output will say so.
//
// Any error go get can output (e.g. "repository does not exist" or "malformed path") are outputted as expected
// and rummage does not touch these kinds of errors
func Get(cmd *cobra.Command, args []string, db *database.Queries, ctx context.Context) {
	flagMap, err := utils.GetBoolFlags(cmd, "update", "dependencies", "debug")
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	flags := getFlags(flagMap)

	if len(args) == 0 && len(flags) > 0 {
		handleFlagsNoPkg(cmd, flags...)
		return
	}

	for _, arg := range args {
		arg = strings.ToLower(arg)
		if strings.Count(arg, "/") >= 2 {
			getAddedItem(cmd, db, ctx, arg, flags...)
			continue
		}
		getHighestScore(cmd, db, ctx, arg, flags...)
	}
}
