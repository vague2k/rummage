package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

func prompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string
	var err error

	var tries int
	for {
		if tries >= 3 {
			fmt.Println("max tries reached for this attempted")
			return false
		}
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, err = r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s = strings.TrimSpace(s)
		if s == "" {
			tries++
		}
		s = strings.ToLower(s)
		if s != "y" && s != "yes" && s != "n" && s != "no" {
			tries++
		}
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func Remove(cmd *cobra.Command, args []string, db *database.Queries, ctx context.Context) {
	flagDeleteAll, err := cmd.Flags().GetBool("delete-all")
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}

	if flagDeleteAll {
		confirmed := prompt("are you sure you want to delete all items?", false)
		if !confirmed {
			return
		}
		err := db.DeleteAllItem(ctx)
		if err != nil {
			cmd.PrintErrf("%s\n", err)
			return
		}
		cmd.Printf("deleted all items from the database\n")
		return
	}

	for _, arg := range args {
		affectedRows, err := db.DeleteItem(ctx, arg)
		if err != nil {
			cmd.PrintErrf("%s\n", arg)
			continue
		} else if affectedRows == 0 {
			cmd.PrintErrf("can't delete item with entry %s it does not exist\n", arg)
			continue
		}
		cmd.Printf("deleted %s from the database\n", arg)
	}
}
