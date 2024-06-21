package database

import (
	"database/sql"
	"errors"
	"fmt"
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

// Creates a predefined schema in the passed db.
func createTable(db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            entry TEXT,
            score FLOAT,
            lastAccessed INTEGER
        )
    `)
	if err != nil {
		msg := fmt.Sprintf("Could not create 'items' table in rummage db: \n%s", err)
		logger.Fatal(msg)
	}
}
