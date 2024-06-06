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

	"github.com/vague2k/rummage/pkg/utils"
)

type RummageDB struct {
	Dir      string // The db's parent directory
	FilePath string // The db.rum file path
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
	dbFile := filepath.Join(dir, "db.rum")

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		msg := fmt.Sprintf("Could not create db dir: \n%s", err)
		return nil, errors.New(msg)
	}

	instance := &RummageDB{
		Dir:      dir,
		FilePath: dbFile,
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
	item := createDBItem(entry, 0.50)

	file, err := os.OpenFile(db.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		msg := fmt.Sprintf("Could not open file path %s for writing: \n%s", db.FilePath, err)
		return errors.New(msg)
	}
	defer file.Close()

	_, err = file.Write(item)
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
	now := utils.Epoch()

	for scanner.Scan() {
		text := scanner.Text()
		entryFromDB := strings.Split(text, "\x00\x00")
		if entry == entryFromDB[0] && fmt.Sprintf("%d", now) != entryFromDB[2] {

			last_accessed, _ := strconv.ParseInt(entryFromDB[2], 10, 64)
			score, _ := strconv.ParseFloat(entryFromDB[1], 0)

			item := &RummageDBItem{
				Entry:        entryFromDB[0],
				Score:        score,
				LastAccessed: last_accessed,
			}
			return true, item
		}
	}

	return false, nil
}

func createDBItem(entry string, defaultScore float64) []byte {
	createdNow := utils.Epoch()
	item := fmt.Sprintf("%s\x00\x00%f\x00\x00%d\n", entry, defaultScore, createdNow)
	b := []byte(item)
	return b
}
