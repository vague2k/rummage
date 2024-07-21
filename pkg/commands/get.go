package commands

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/vague2k/rummage/logger"
	"github.com/vague2k/rummage/pkg/database"
)

var log = logger.New()

// Attempts to call "go get" against an arg and takes into account any flags that "go get" may accept.
func AttemptGoGet(arg string, flags ...string) error {
	var cmd *exec.Cmd
	if flags == nil {
		cmd = exec.Command("go", "get", arg)
	} else {
		argWithFlags := []string{"get"}
		for _, f := range flags {
			argWithFlags = append(argWithFlags, f)
		}
		cmd = exec.Command("go", argWithFlags...)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		// if there's an error, []byte returned from cmd.CombinedOutput() will contain the error
		msg := fmt.Sprintf("%v", string(b))
		return errors.New(msg)
	}

	log.Info("Got the following packages...\n", string(b))

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
		log.Fatal("Did not incrment item's score due to error: \n", err)
	}

	return item
}

// Attempts to call "go get" against an item and if it does not exist in the db, adds it.
func GoGetAddedItem(db *database.RummageDB, arg string, flags ...string) *database.RummageDBItem {
	if err := AttemptGoGet(arg, flags...); err != nil {
		log.Err(err)
	}

	added, err := db.AddItem(arg)
	if err != nil {
		log.Fatal(err) // if db can't be accessed
	}

	UpdateRecency(db, added)
	return added
}

// Attempts to call "go get" against an item with the highest score where the substr matches any items found in the db.
func GoGetHighestScore(db *database.RummageDB, substr string, flags ...string) {
	found, exists := db.EntryWithHighestScore(substr)
	if !exists {
		log.Warn("No entry was found that could match your argument.")
		return
	}

	err := AttemptGoGet(found.Entry, flags...)
	if err != nil {
		log.Err(err)
	}

	UpdateRecency(db, found)
}
