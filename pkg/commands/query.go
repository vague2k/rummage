package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/pkg/commands/utils"
	"github.com/vague2k/rummage/pkg/database"
)

// The Query command lets a user query the database against an arg to return the entry.
// By default with no flags the query will search based on highest score, if the "--exact" flag is true,
// it will simply do a strings.contains search.
//
// You can view other info like the entry's score and last-accessed field by setting the respective flags,
// "--score", "--last-accessed"
func Query(cmd *cobra.Command, arg string, db database.RummageDbInterface) {
	flagMap, err := utils.GetBoolFlags(cmd, "exact", "score", "last-accessed")
	if err != nil {
		cmd.PrintErrf("could not register a flag for 'get', cobra gave this error: \n%s\n", err)
	}

	var found *database.RummageItem
	var s strings.Builder

	arg = strings.ToLower(arg)

	// the method by which we query has to be either exclusively "exact" or by score, thus the switch statement
	// I could make searching by score it's own flag and mark both as mutually exclusive, but i think this for now
	// gets the point across well enough
	switch true {
	case flagMap["exact"]:
		item, err := db.FindExactMatch(arg)
		if err != nil {
			cmd.PrintErrf("%s\n", err)
			return
		}
		found = item
		s.WriteString(fmt.Sprintf("%s ", found.Entry))
	case !flagMap["exact"]:
		item, err := db.EntryWithHighestScore(arg)
		if err != nil {
			cmd.PrintErrf("%s\n", err)
			return
		}
		found = item
		s.WriteString(fmt.Sprintf("%s ", found.Entry))
	}

	if flagMap["score"] {
		s.WriteString(fmt.Sprintf("%.2f ", found.Score))
	}

	if flagMap["last-accessed"] {
		s.WriteString(fmt.Sprintf("%d ", found.LastAccessed))
	}

	cmd.Printf("%s\n", s.String())
}
