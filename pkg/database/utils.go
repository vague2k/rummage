package database

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// Gets the user's data dir
//
// ".local/share" on unix, %LOCALAPPDATA% on windows
func dataDir() (string, error) {
	// TODO: account for XDG env vars
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
