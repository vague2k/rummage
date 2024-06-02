package db

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RummageDB struct {
	Dir      string
	FilePath string
}

// Accesses the rummage db, returning a pointer to the db instance.
//
// Access() also makes sure the directory exists, but does not write anything to it's children.
func Access() (*RummageDB, error) {
	dir := "/Users/albert/.local/share/rummage"

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		msg := fmt.Sprint("Could not create db dir: \n", err)
		return nil, errors.New(msg)
	}

	instance := &RummageDB{
		Dir:      dir,
		FilePath: filepath.Join(dir, "db.rum"),
	}

	return instance, nil
}

// Adds an item to the db.
// If the db.FilePath does not exist, it will be created.
//
// If the item's entry already exists, AddItem() does nothing.
func (db *RummageDB) AddItem(entry string, score int) error {
	if db.EntryExists(entry) {
		return nil
	}
	// item uses a double null byte to distinguish between the entry and it's score,
	item := entry + "\x00\x00" + fmt.Sprintf("%d", score) + "\n"
	itemBytes := []byte(item)

	file, err := os.OpenFile(db.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		msg := fmt.Sprintf("Could not open file path %s for writing: \n%s", db.FilePath, err)
		return errors.New(msg)
	}
	defer file.Close()

	_, err = file.Write(itemBytes)
	if err != nil {
		msg := fmt.Sprintf("Issue occured when writing item to file path %s: \n%s", db.FilePath, err)
		return errors.New(msg)
	}

	return nil
}

// Checks if a specific db.Item.Entry exists in the db.
// If the entry does not exist, returns false.
//
// If the db cannot be read, an error will be propagated.
func (db *RummageDB) EntryExists(entry string) bool {
	file, err := os.Open(db.FilePath)
	if err != nil {
		log.Fatalf("Could not open file path %s for reading: \n%s", db.FilePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		entryFromDB := strings.Split(text, "\x00\x00")
		if entry == entryFromDB[0] {
			return true
		}
	}

	return false
}
