package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/internal"
)

var logger = internal.NewLogger(nil, os.Stdout)

func RegisterBoolFlags(cmd *cobra.Command, names ...string) map[string]bool {
	flags := make(map[string]bool)

	for _, name := range names {
		flagVal, err := cmd.Flags().GetBool(name)
		if err != nil {
			logger.Err("Could not register ", name, " flag.")
		}
		flags[name] = flagVal
	}

	return flags
}
