package db

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vague2k/rummage/pkg/utils"
)

type RummageDB struct {
	Dir      string
	FilePath string
}

// Accesses the rummage db, returning a pointer to the db instance.
//
// Access() also makes sure the directory exists, but does not write anything to it's children.
func Access() (*RummageDB, error) {
	dataDir, err := utils.DataDir()
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	dir := filepath.Join(dataDir, "rummage")

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		msg := fmt.Sprintf("Could not create db dir: \n%s", err)
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
func (db *RummageDB) AddItem(entry string) error {
	if exists, _ := db.EntryExists(entry); exists {
		return nil
	}
	now := time.Now().Unix()
	// item uses a double null byte to distinguish between the entry and it's score,
	item := fmt.Sprintf("%s\x00\x00%d\x00\x00%d\n", entry, 1, now)
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

// Checks if a specific db.Item exists in the db.
//
// If the item does exist, returns true and a pointer to a RummageDBItem.
// If the item does not exist, returns false and nil.
//
// If the db cannot be read, an error will be propagated.
func (db *RummageDB) EntryExists(entry string) (bool, *RummageDBItem) {
	file, err := os.Open(db.FilePath)
	if err != nil {
		log.Fatalf("Could not open file path %s for reading: \n%s", db.FilePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	now := time.Now().Unix()

	for scanner.Scan() {
		text := scanner.Text()
		entryFromDB := strings.Split(text, "\x00\x00")
		if entry == entryFromDB[0] && fmt.Sprintf("%d", now) != entryFromDB[2] {

			last_accessed, _ := strconv.ParseInt(entryFromDB[2], 10, 64)
			score, _ := strconv.ParseInt(entryFromDB[1], 10, 0)

			item := &RummageDBItem{
				Entry:        entryFromDB[0],
				Score:        score,
				LastAccessed: time.Unix(last_accessed, 0),
			}
			return true, item
		}
	}

	return false, nil
}
