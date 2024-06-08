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

// Adds an item to the db and returns a pointer to the item that was just added.
// Newly added items are given a default score of 1.0.
//
// If the db.FilePath does not exist, it will be created.
//
// If the item's entry already exists, AddItem() returns the item
func (db *RummageDB) AddItem(entry string) (*RummageDBItem, error) {
	// make sure the file exists first before checking if the entry exists,
	// to make sure "db.rum" is created if that file does not exist
	file, err := os.OpenFile(db.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		msg := fmt.Sprintf("Could not open file path %s for writing: \n%s", db.FilePath, err)
		return nil, errors.New(msg)
	}
	defer file.Close()

	if exists, item := db.EntryExists(entry); exists {
		return item, nil
	}

	defaultScore := 1.0
	item := createDBItem(entry, defaultScore, true)

	_, err = file.Write(item)
	if err != nil {
		msg := fmt.Sprintf("Issue occured when writing item to file path %s: \n%s", db.FilePath, err)
		return nil, errors.New(msg)
	}

	return &RummageDBItem{
		Entry:        entry,
		Score:        defaultScore,
		LastAccessed: time.Now().Unix(),
	}, nil
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

	for scanner.Scan() {
		text := scanner.Text()
		entryFromDB := strings.Split(text, "\x00\x00")
		if entry == entryFromDB[0] {

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

// Updates an item in the db if the entry can be found. An entry not being found is treated as an error.
//
// If the db file can't be read, written to, or if the entry can't be found, an error will also be returned.
func (db *RummageDB) UpdateItem(entry string, updated *RummageDBItem) (*RummageDBItem, error) {
	contents, err := os.ReadFile(db.FilePath)
	if err != nil {
		msg := fmt.Sprintf("Could not open file path %s for reading: \n%s", db.FilePath, err)
		return nil, errors.New(msg)
	}

	if exists, _ := db.EntryExists(entry); !exists {
		msg := fmt.Sprintf("The entry, %s could not be found", entry)
		return nil, errors.New(msg)
	}

	lines := strings.Split(string(contents), "\n")
	for i, line := range lines {
		item := strings.Split(line, "\x00\x00")
		if item[0] == entry {
			lines[i] = string(createDBItem(updated.Entry, updated.Score, false))
		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(db.FilePath, []byte(output), 0644)
	if err != nil {
		log.Fatalf("Could not open update (write) item at path %s: \n%s", db.FilePath, err)
	}

	updatedItem := &RummageDBItem{
		Entry:        updated.Entry,
		Score:        updated.Score,
		LastAccessed: time.Now().Unix(),
	}

	return updatedItem, nil
}

func createDBItem(entry string, defaultScore float64, newline bool) []byte {
	var item string
	createdNow := time.Now().Unix()
	if newline {
		item = fmt.Sprintf("%s\x00\x00%f\x00\x00%d\n", entry, defaultScore, createdNow)
	} else {
		item = fmt.Sprintf("%s\x00\x00%f\x00\x00%d", entry, defaultScore, createdNow)
	}
	b := []byte(item)
	return b
}
