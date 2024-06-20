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
func attemptGoGet(arg string) error {
	cmd := exec.Command("go", "get", arg)
	_, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("Could not combined output from 'go get' cmd: \n%s", err)
		return errors.New(msg)
	}

	// output := string(b)
	// log.Print(output)

	return nil
}

// Attempts to call "go get" against an item and if it does not exist in the db, adds it.
func GoGetAddedItem(db *database.RummageDB, arg string) *database.RummageDBItem {
	if err := attemptGoGet(arg); err != nil {
		log.Print(err)
		return nil
	}

	added, err := db.AddItem(arg)
	if err != nil {
		log.Fatal(err) // if db can't be accessed
	}

	increment := database.RummageDBItem{
		Entry:        added.Entry,
		Score:        added.RecalculateScore(),
		LastAccessed: time.Now().Unix(),
	}

	_, err = db.UpdateItem(added.Entry, &increment)
	return added
}

// Attempts to call "go get" against an item with the highest score where the substr matches any items found in the db.
func GoGetHighestScore(db *database.RummageDB, substr string) {
	found, exists := db.EntryWithHighestScore(substr)
	if !exists {
		GoGetAddedItem(db, substr) // will error early if substr so happens to not be a go package
		return
	}

	fmt.Printf("got here 0")
	err := attemptGoGet(found.Entry)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got here 1")

	increment := database.RummageDBItem{
		Entry:        found.Entry,
		Score:        found.RecalculateScore(),
		LastAccessed: time.Now().Unix(),
	}

	fmt.Printf("got here 2")
	_, err = db.UpdateItem(found.Entry, &increment)
	if err != nil {
		log.Printf("Did not incrment item's score due to error: \n%s", err)
	}
	fmt.Printf("got here 3")
}
