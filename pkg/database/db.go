package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vague2k/rummage/pkg/config"
	"github.com/vague2k/rummage/utils"
)

// why does this interface exist even though there's only 1 implementation
//
// because it makes my life writing tests very easy.
type RummageDbInterface interface {
	Version() string
	AddItem(entry string) (*RummageItem, error)
	AddMultiItems(entries ...string) ([]*RummageItem, int, error)
	Close()
	DeleteAllItems() error
	DeleteItem(entry string) (*RummageItem, error)
	EntryWithHighestScore(substr string) (*RummageItem, error)
	FindTopNMatches(substr string, n int) ([]*RummageItem, error)
	ListItems() ([]*RummageItem, error)
	SelectItem(entry string) (*RummageItem, error)
	UpdateItem(entry string, score float64, lastAccessed int64) (*RummageItem, error)
}

// A database wrapper for Rummage that pertains to rummage's actions
//
// This struct and it's methods implements RummageDbInterface
type RummageDb struct {
	Sqlite   *sql.DB // Pointer to the underlying sqlite database
	Dir      string  // The parent directory of the database
	FilePath string  // the database path
	version  string
}

// Initializes the rummage db, returning a pointer to the db instance.
//
// Init() also makes sure the "items" table exists.
func Init(path string) (*RummageDb, error) {
	var ver string
	if path == "" {
		dataDir := utils.UserDataDir()
		path = dataDir

		conf := config.SetVersions()
		ver = conf.Rummage.DbApiVersion
	}

	var dir string
	var dbFile string
	if path == ":memory:" {
		dbFile = ":memory:"
		ver = ""
	} else {
		dir = filepath.Join(path, "rummage")
		dbFile = filepath.Join(dir, "rummage.db")
		err := os.MkdirAll(dir, 0o777)
		if err != nil {
			return nil, fmt.Errorf("could not create db dir: \n%s", err)
		}
	}

	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not init rummage db: \n%s", err)
	}

	_, err = database.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            entry TEXT,
            score FLOAT,
            lastAccessed INTEGER
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("could not create 'items' table in rummage db: \n%s", err)
	}

	instance := &RummageDb{
		Sqlite:   database,
		Dir:      dir,
		FilePath: dbFile,
		version:  ver,
	}

	return instance, nil
}

func (r *RummageDb) Version() string {
	return r.version
}

func (r *RummageDb) Close() {
	r.Sqlite.Close()
}

// Adds an item to the db and returns a pointer to the item that was just added.
// Newly added items are given a default score of 1.0.
//
// If the item's entry already exists, AddItem() returns the item
func (r *RummageDb) AddItem(entry string) (*RummageItem, error) {
	var item *RummageItem

	if item, err := r.SelectItem(entry); err == nil {
		return item, nil
	}

	regex := regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z0-9-!]+(/[a-zA-Z0-9-_\.!]+)+(/[vV]\d+)?$`)
	if !regex.MatchString(entry) {
		return nil, fmt.Errorf("the item attempted to be added to the database does not resemble a go package")
	}

	_, err := r.Sqlite.Exec(`
        INSERT INTO items (entry, score, lastAccessed) 
        VALUES (?, ?, ?)`,
		entry, 1.0, time.Now().Unix(),
	)
	if err != nil {
		return nil, fmt.Errorf("issue occured when adding item to db: \n%s", err)
	}

	item = &RummageItem{
		Entry:        entry,
		Score:        1.0,
		LastAccessed: time.Now().Unix(),
	}
	return item, nil
}

// Selects a specific item in the db by it's entry, and returns a pointer to it.
//
// If the item does not exist, it is treated as an error
func (r *RummageDb) SelectItem(entry string) (*RummageItem, error) {
	var item RummageItem

	row := r.Sqlite.QueryRow(`
        SELECT * FROM items WHERE entry = ? 
        LIMIT 1
        `,
		entry,
	)

	err := row.Scan(&item.Entry, &item.Score, &item.LastAccessed)
	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("the item with entry %s does not exist", entry)
	}

	return &item, nil
}

// Updates an item in the db if the entry can be found.
func (r *RummageDb) UpdateItem(entry string, score float64, lastAccessed int64) (*RummageItem, error) {
	if _, err := r.SelectItem(entry); err != nil {
		return nil, fmt.Errorf("the item with entry %s is attempted to be updated but does not exist", entry)
	}

	_, err := r.Sqlite.Exec(`
        UPDATE items
        SET score = ?, lastAccessed = ?
        WHERE entry = ?
        `,
		score,
		lastAccessed,
		entry,
	)
	if err != nil {
		return nil, fmt.Errorf("issue occured while trying to update db item: \n%s", err)
	}

	updatedItem := &RummageItem{
		Entry:        entry,
		Score:        score,
		LastAccessed: lastAccessed,
	}

	return updatedItem, nil
}

// List all items in the db, returning []RummageDBItem
func (r *RummageDb) ListItems() ([]*RummageItem, error) {
	rows, err := r.Sqlite.Query("SELECT * FROM items")
	if err != nil {
		return nil, fmt.Errorf("could not get items from db: \n%s", err)
	}

	var items []*RummageItem
	for rows.Next() {
		var entry string
		var score float64
		var lastAccessed int64

		err := rows.Scan(&entry, &score, &lastAccessed)
		if err != nil {
			return nil, fmt.Errorf("issue occured trying to scan a db item: \n%s", err)
		}

		nextItem := &RummageItem{
			Entry:        entry,
			Score:        score,
			LastAccessed: lastAccessed,
		}

		items = append(items, nextItem)
	}

	return items, nil
}

// Searches the db for an entry that matches a substring "substr",
// if multiple matches are found, the item with the highest scores is returned
//
// No matches is treated as an error
func (r *RummageDb) EntryWithHighestScore(substr string) (*RummageItem, error) {
	items, err := r.ListItems()
	if err != nil {
		log.Fatal(err)
	}

	if len(items) == 1 {
		return items[0], nil
	}

	var curr float64
	var highestMatch *RummageItem
	for _, item := range items {
		if strings.Contains(item.Entry, substr) {
			if item.Score > curr {
				highestMatch = item
				curr = item.Score
			}
		}
	}

	if highestMatch == nil {
		return nil, fmt.Errorf("no match found with the given arguement %s", substr)
	}

	return highestMatch, nil
}

// Finds top n matches in the database based on a substr, sorted in descending by score
//
// No matches is treated as an error
func (r *RummageDb) FindTopNMatches(substr string, n int) ([]*RummageItem, error) {
	rows, err := r.Sqlite.Query(
		"SELECT * FROM items WHERE entry LIKE ? ORDER BY score DESC LIMIT ?",
		"%"+substr+"%", n,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get items from db: \n%s", err)
	}

	items := make([]*RummageItem, 0, n)
	for rows.Next() {
		nextItem := &RummageItem{}

		err := rows.Scan(&nextItem.Entry, &nextItem.Score, &nextItem.LastAccessed)
		if err != nil {
			return nil, fmt.Errorf("issue occured trying to scan a db item: \n%s", err)
		}

		items = append(items, nextItem)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no match found with the given arguement %s", substr)
	}

	return items, nil
}

// Adds multiple items to the db and returns []*RummageDBItem along with the number of items added
func (r *RummageDb) AddMultiItems(entries ...string) ([]*RummageItem, int, error) {
	var slice []*RummageItem
	var itemsAdded int

	for _, entry := range entries {
		// if the entry attempted to be added already exists, skip this iteration
		if _, err := r.SelectItem(entry); err == nil {
			continue
		}

		item, err := r.AddItem(entry)
		if err != nil {
			return slice, itemsAdded, fmt.Errorf("issue occured when adding item %s to the db", entry)
		}

		slice = append(slice, item)
		itemsAdded++
	}

	return slice, itemsAdded, nil
}

// Deletes an item from the database, and returns a pointer to the item that was just deleted.
func (r *RummageDb) DeleteItem(entry string) (*RummageItem, error) {
	item, err := r.SelectItem(entry)
	if err != nil {
		return nil, fmt.Errorf("can't delete item with entry %s it does not exist", entry)
	}
	_, err = r.Sqlite.Exec(`
        DELETE FROM items
        WHERE entry = ?
        `,
		entry,
	)
	if err != nil {
		return nil, fmt.Errorf("issue occured while trying to delete db item: \n%s", err)
	}

	return item, nil
}

// Deletes all items from the database, should be used with a yes/no prompt of some sort.
func (r *RummageDb) DeleteAllItems() error {
	_, err := r.Sqlite.Exec("DELETE FROM items")
	if err != nil {
		return fmt.Errorf("issue occured when trying to delete all items in the db: \n%s", err)
	}
	return nil
}
