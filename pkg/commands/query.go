package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

// The Query command lets a user query the database against an arg to output a list of
// entries (10 by default) sorted by score.
//
// the quantity of matches in the output can be changed with the "--quantity" flag
func Query(cmd *cobra.Command, arg string, db *database.Queries, ctx context.Context) {
	quantityFlag, err := cmd.Flags().GetInt("quantity")
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}

	arg = strings.ToLower(arg)
	items, err := db.FindTopNMatches(ctx, database.FindTopNMatchesParams{
		Entry: "%" + arg + "%",
		Limit: int64(quantityFlag),
	})
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	} else if len(items) == 0 {
		cmd.PrintErrf("%s %s\n", "no match found with the given arguement", arg)
		return
	}
	var s strings.Builder

	// formatting
	var entryMaxLen, lastAccessedMaxLen, scoreMaxLen int
	for _, item := range items {
		entryLen := len(item.Entry)
		scoreLen := len(fmt.Sprintf("%.4f", item.Score))
		lastAccessedLen := len(fmt.Sprintf("%d", item.Lastaccessed))

		if entryLen > entryMaxLen {
			entryMaxLen = entryLen
		}
		if scoreLen > scoreMaxLen {
			scoreMaxLen = scoreLen
		}
		if lastAccessedLen > lastAccessedMaxLen {
			lastAccessedMaxLen = lastAccessedLen
		}
	}

	// Formatting output with proper padding
	for _, item := range items {
		s.WriteString(fmt.Sprintf(
			"%-*d : %-*.*f : %-*s\n",
			lastAccessedMaxLen, item.Lastaccessed,
			scoreMaxLen, 4, item.Score,
			entryMaxLen, item.Entry,
		))
	}

	cmd.Printf("%s\n", s.String())
}
