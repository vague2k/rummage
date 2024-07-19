package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vague2k/rummage/logger"
)

var log = logger.New()

// Creates a "items" table if it does not exist in the Rummage DB.
func CreateDBTable(db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            entry TEXT,
            score FLOAT,
            lastAccessed INTEGER
        )
    `)
	if err != nil {
		msg := fmt.Sprintf("Could not create 'items' table in rummage db: \n%s", err)
		log.Fatal(msg)
	}
}

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
		log.Fatal("\n", err)
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
		log.Fatal("Could not get the home directory")
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

func RegisterBoolFlags(cmd *cobra.Command, names ...string) map[string]bool {
	flags := make(map[string]bool)

	for _, name := range names {
		flagVal, err := cmd.Flags().GetBool(name)
		if err != nil {
			log.Err("Could not register ", name, " flag.")
		}
		flags[name] = flagVal
	}

	return flags
}
