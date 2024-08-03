package utils

import (
	"github.com/spf13/cobra"
)

// Registers a map of key value pairs where the key are flag names (string)
// and the values are wether the flag is on/off (bool)
func GetBoolFlags(cmd *cobra.Command, names ...string) (map[string]bool, error) {
	flags := make(map[string]bool)
	for _, name := range names {
		v, err := cmd.Flags().GetBool(name)
		if err != nil {
			return nil, err
		}
		flags[name] = v
	}

	return flags, nil
}
