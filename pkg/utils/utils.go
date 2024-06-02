package utils

import (
	"errors"
	"os"
)

// Returns the array of args excluding the command name.
//
// An empty array (no args) is treated as an error.
func ParseArgs() ([]string, error) {
	args := os.Args

	if len(args) <= 1 {
		return nil, errors.New("There were no args parsed, array is empty")
	}

	return args[1:], nil
}
