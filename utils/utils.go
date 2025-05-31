package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

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

func ResemblesGoPackage(entry string) error {
	regex := regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z0-9-!]+(/[a-zA-Z0-9-_\.!]+)+(/[vV]\d+)?$`)
	if !regex.MatchString(entry) {
		return fmt.Errorf("the item attempted to be added to the database does not resemble a go package")
	}
	return nil
}
