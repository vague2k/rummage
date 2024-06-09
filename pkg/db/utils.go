package db

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Creates an db item object returned as []byte
//
// If newline is true, add's a newline to the end of the item string. Useful for updating a item, which would use false.
func createDBItem(entry string, score float64, newline bool) []byte {
	createdNow := time.Now().Unix()
	item := fmt.Sprintf("%s\x00\x00%f\x00\x00%d", entry, score, createdNow)

	if newline {
		item += "\n"
	}

	b := []byte(item)
	return b
}

// Gets the user's data dir
//
// ".local" on unix, %LOCALAPPDATA% on windows
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

// Takes a file path and open's it, allowing you to scan over it performing some operation using *bufio.scanner.
//
// The file opened in RONLY and is defered to close after the callback's execution is complete.
func scanOverFile(path string, callback func(scanner *bufio.Scanner)) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not open file path %s for reading: \n%s", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		callback(scanner)
	}
}
