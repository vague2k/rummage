package pkg

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/vague2k/rummage/pkg/database"
)

// Attempts to call "go get" against an arg.
func AttemptGoGet(arg string) error {
	cmd := exec.Command("go", "get", arg)
	b, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("Could not combined output from 'go get' cmd: \n%s", err)
		return errors.New(msg)
	}

	output := string(b)
	fmt.Printf("%s", output)

	return nil
}

// Increments an item's score and updates the LastAccessed field to time.Now().Unix().
func UpdateRecency(db *database.RummageDB, item *database.RummageDBItem) *database.RummageDBItem {
	recency := &database.RummageDBItem{
		Entry:        item.Entry,
		Score:        item.RecalculateScore(),
		LastAccessed: time.Now().Unix(),
	}
	item, err := db.UpdateItem(item.Entry, recency)
	if err != nil {
		log.Fatalf("Did not incrment item's score due to error: \n%s", err)
	}

	return item
}

// Attempts to call "go get" against an item and if it does not exist in the db, adds it.
func GoGetAddedItem(db *database.RummageDB, arg string) *database.RummageDBItem {
	if err := AttemptGoGet(arg); err != nil {
		log.Print(err)
		return nil
	}

	added, err := db.AddItem(arg)
	if err != nil {
		log.Fatal(err) // if db can't be accessed
	}

	UpdateRecency(db, added)
	return added
}

// Attempts to call "go get" against an item with the highest score where the substr matches any items found in the db.
func GoGetHighestScore(db *database.RummageDB, substr string) {
	found, exists := db.EntryWithHighestScore(substr)
	if !exists {
		fmt.Printf("%s", "No entry was found that could match your argument.")
		return
	}

	err := AttemptGoGet(found.Entry)
	if err != nil {
		log.Fatal(err)
	}

	UpdateRecency(db, found)
}
