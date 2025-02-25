package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/database"
)

// The Query command lets a user query the database against an arg to return the entry.
// By default with no flags the query will search based on highest score, if the "--exact" flag is true,
// it will simply do a strings.contains search.
//
// You can view other info like the entry's score and last-accessed field by setting the respective flags,
// "--score", "--last-accessed"
func Query(cmd *cobra.Command, arg string, db database.RummageDbInterface) {
	quantityFlag, err := cmd.Flags().GetInt("quantity")
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}

	arg = strings.ToLower(arg)
	items, err := db.FindTopNMatches(arg, quantityFlag)
	if err != nil {
		cmd.PrintErrf("%s\n", err)
		return
	}
	var s strings.Builder

	// formatting
	var entryMaxLen, lastAccessedMaxLen, scoreMaxLen int
	for _, item := range items {
		entryLen := len(item.Entry)
		scoreLen := len(fmt.Sprintf("%.4f", item.Score))
		lastAccessedLen := len(fmt.Sprintf("%d", item.LastAccessed))

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
			lastAccessedMaxLen, item.LastAccessed,
			scoreMaxLen, 4, item.Score,
			entryMaxLen, item.Entry,
		))
	}

	cmd.Printf("%s\n", s.String())
}
