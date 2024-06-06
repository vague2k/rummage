package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"time"
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

func DataDir() (string, error) {
	var dataDir string

	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Could not get the home directory")
	}

	switch runtime.GOOS {
	case "windows":
		dataDir = os.Getenv("LOCALAPPDATA")
		if dataDir == "" {
			dataDir = filepath.Join(home, "AppData", "Local")
		}
	default:
		dataDir = filepath.Join(home, ".local", "share")
	}

	return dataDir, nil
}

// Wrapper for time.Now().Unix(), just in case the specific unix time (currently returns ms), needs to be changed on the fly.
//
// this util func may or may not stay here / be used in prod.
func Epoch() int64 {
	return time.Now().Unix()
}
