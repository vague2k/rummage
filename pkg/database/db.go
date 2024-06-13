package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

	// open the db
	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		msg := fmt.Sprintf("Could not init rummage db: \n%s", err)
		return nil, errors.New(msg)
	}

	// create the items table if it doesn't exist
	_, err = database.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            entry TEXT,
            score FLOAT,
            lastAccessed INTEGER
        )
    `)
	if err != nil {
		msg := fmt.Sprintf("Could not create 'items' table in rummage db: \n%s", err)
		return nil, errors.New(msg)
	}

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

	newEntry := updated.Entry
	newScore := updated.Score
	newLastAccessed := time.Now().Unix()

	_, err := r.DB.Exec(`
        UPDATE items
        SET score = ?, lastAccessed = ?
        WHERE entry = ?
        `,
		newScore,
		newLastAccessed,
		entry,
	)
	if err != nil {
		msg := fmt.Sprintf("Issue updating db item: \n%s", err)
		return nil, errors.New(msg)
	}

	updatedItem := &RummageDBItem{
		Entry:        newEntry,
		Score:        newScore,
		LastAccessed: newLastAccessed,
	}

	return updatedItem, nil
}