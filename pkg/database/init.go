package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

func Init(path string) (*Queries, error) {
	if path == "" {
		dataDir := userDataDir()
		path = dataDir
	}

	var dir string
	var dbFile string
	if path == ":memory:" {
		dbFile = ":memory:"
	} else {
		dir = filepath.Join(path, "rummage")
		dbFile = filepath.Join(dir, "rummage.db")
		err := os.MkdirAll(dir, 0o777)
		if err != nil {
			return nil, fmt.Errorf("could not create db dir: \n%s", err)
		}
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not init rummage db: \n%s", err)
	}

	if _, err := db.ExecContext(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("could not create 'items' table in rummage db: \n%s", err)
	}

	return New(db), nil
}

// Gets the user's $XDG_DATA_HOME dir.
//
// Fallsback to the default data dir if the env var does not exist.
func userDataDir() string {
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
