package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/vague2k/rummage/internal"
)

var logger = internal.NewLogger(nil, os.Stdout)

// Parses "/" from a string.
//
// If a match is found, true and the length of matches is returned.
func ParseForwardSlash(s string) (bool, int) {
	m := regexp.MustCompile("/").FindAllString(s, -1)
	return len(m) != 0, len(m)
}

// Gets the user's $GOPATH.
//
// Gets "user-home-dir/go" if $GOPATH does not exist.
func UserGoPath() string {
	if path := os.Getenv("GOPATH"); path != "" {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal("\n", err)
	}

	return filepath.Join(home, "go")
}

// Gets the user's $XDG_DATA_HOME dir.
//
// Fallsback to the default data dir if the env var does not exist.
func UserDataDir() string {
	var dataDir string

	if dataDir = os.Getenv("XDG_DATA_HOME"); dataDir != "" {
		return dataDir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal("Could not get the home directory")
	}

	switch runtime.GOOS {
	case "windows":
		if dataDir = os.Getenv("LOCALAPPDATA"); dataDir != "" {
			dataDir = filepath.Join(home, "AppData", "Local")
		}
	default:
		dataDir = filepath.Join(home, ".local", "share")
	}

	return dataDir
}
