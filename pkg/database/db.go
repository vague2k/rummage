package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// A database wrapper for Rummage that pertains to rummage's actions
type RummageDB struct {
	DB       *sql.DB // Pointer to the underlying sqlite database
	Dir      string  // The parent directory of the database
	FilePath string  // the database path
}

// Inits the rummage db, returning a pointer to the db instance.
//
// Init() also makes sure the "items" table exists.
func Init(path string) (*RummageDB, error) {
	if path == "" {
		dataDir, err := dataDir()
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
		path = dataDir
	}

	// make sure the path to the db exists
	// TODO: may need to rethink this if we want to use a user defined config file
	dir := filepath.Join(path, "rummage")
	dbFile := filepath.Join(dir, "rummage.db")

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		msg := fmt.Sprintf("Could not create db dir: \n%s", err)
		return nil, errors.New(msg)
	}

	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		msg := fmt.Sprintf("Could not init rummage db: \n%s", err)
		return nil, errors.New(msg)
	}
	createTable(database) // create the items table if it doesn't exist

	instance := &RummageDB{
		DB:       database,
		Dir:      dir,
		FilePath: dbFile,
	}

	return instance, nil
}

// Adds an item to the db and returns a pointer to the item that was just added.
// Newly added items are given a default score of 1.0.
//
// If the item's entry already exists, AddItem() returns the item
func (r *RummageDB) AddItem(entry string) (*RummageDBItem, error) {
	var item RummageDBItem

	if item, exists := r.SelectItem(entry); exists {
		return item, nil
	}

	_, err := r.DB.Exec(`
        INSERT INTO items (entry, score, lastAccessed) 
        VALUES (?, ?, ?)`,
		entry, 1.0, time.Now().Unix(),
	)
	if err != nil {
		msg := fmt.Sprintf("Issue occured when adding item to db: \n%s", err)
		return nil, errors.New(msg)
	}

	item = RummageDBItem{
		Entry:        entry,
		Score:        1.0,
		LastAccessed: time.Now().Unix(),
	}
	return &item, nil
}

// Selects a specific item in the db by it's entry, and returns a pointer to it.
//
// If the item does not exist, returns nil, false
func (r *RummageDB) SelectItem(entry string) (*RummageDBItem, bool) {
	var item RummageDBItem

	row := r.DB.QueryRow(`
        SELECT * FROM items WHERE entry = ? 
        LIMIT 1
        `,
		entry,
	)

	err := row.Scan(&item.Entry, &item.Score, &item.LastAccessed)
	if err != nil && err == sql.ErrNoRows {
		return nil, false
	}

	return &item, true
}

// Updates an item in the db if the entry can be found.
//
// An entry not being found is treated as an error.
func (r *RummageDB) UpdateItem(entry string, updated *RummageDBItem) (*RummageDBItem, error) {
	if _, exists := r.SelectItem(entry); !exists {
		msg := fmt.Sprintf("The entry, %s could not be found", entry)
		return nil, errors.New(msg)
	}

	_, err := r.DB.Exec(`
        UPDATE items
        SET score = ?, lastAccessed = ?
        WHERE entry = ?
        `,
		updated.Score,
		updated.LastAccessed,
		entry,
	)
	if err != nil {
		msg := fmt.Sprintf("Issue updating db item: \n%s", err)
		return nil, errors.New(msg)
	}

	updatedItem := &RummageDBItem{
		Entry:        entry,
		Score:        updated.Score,
		LastAccessed: updated.LastAccessed,
	}

	return updatedItem, nil
}

// List all items in the db, returning []RummageDBItem
//
// An error is thrown if there was an issue scanning a row.
func (r *RummageDB) ListItems() ([]RummageDBItem, error) {
	rows, err := r.DB.Query("SELECT * FROM items")
	if err != nil {
		msg := fmt.Sprintf("Could not get items from db: \n%s", err)
		return nil, errors.New(msg)
	}

	var items []RummageDBItem
	for rows.Next() {
		count := 0 // track which iteration of the loop incase of err when scanning rows
		var entry string
		var score float64
		var lastAccessed int64

		err := rows.Scan(&entry, &score, &lastAccessed)
		if err != nil {
			msg := fmt.Sprintf("Could not get (%d)th iteration of item from db: \n%s", count, err)
			return nil, errors.New(msg)
		}

		nextItem := RummageDBItem{
			Entry:        entry,
			Score:        score,
			LastAccessed: lastAccessed,
		}

		items = append(items, nextItem)
		count++
	}

	return items, nil
}

// Searches the db for an entry that matches a substring "substr",
// if multiple matches are found, the match with the highest score's *RummageDBItem is returned
// If no match was found, false is returned
func (r *RummageDB) EntryWithHighestScore(substr string) (*RummageDBItem, bool) {
	items, err := r.ListItems()
	if err != nil {
		log.Fatal(err)
	}

	var curr float64
	var highestMatch *RummageDBItem

	for _, item := range items {
		if strings.Contains(item.Entry, substr) {
			if item.Score > curr {
				highestMatch = &item
				curr = item.Score
			}
		}
	}

	if highestMatch == nil {
		return nil, false
	}

	return highestMatch, true
}
