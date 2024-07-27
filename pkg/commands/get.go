package commands

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/vague2k/rummage/logger"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/utils"
)

var log = logger.New()

// Attempts to call "go get" against an arg and takes into account any flags that "go get" may accept.
func attemptGoGet(arg string, flags ...string) error {
	var cmd *exec.Cmd
	if flags == nil {
		cmd = exec.Command("go", "get", arg)
	} else {
		getWithFlags := append([]string{"get"}, flags...)
		argWithFlags := append(getWithFlags, arg)
		cmd = exec.Command("go", argWithFlags...)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		// if there's an error, []byte returned from cmd.CombinedOutput() will contain the error
		return errors.New(string(b))
	}

	log.Print(string(b))
	return nil
}

// Increments an item's score and updates the LastAccessed field to time.Now().Unix().
func updateRecency(db *database.RummageDB, item *database.RummageDBItem) error {
	recency := &database.RummageDBItem{
		Entry:        item.Entry,
		Score:        item.RecalculateScore(),
		LastAccessed: time.Now().Unix(),
	}
	_, err := db.UpdateItem(item.Entry, recency)
	if err != nil {
		return err
	}
	return nil
}

// Attempts to call "go get" against an item and if it does not exist in the db, adds it.
func goGetAddedItem(db *database.RummageDB, arg string, flags ...string) error {
	if err := attemptGoGet(arg, flags...); err != nil {
		return err
	}

	added, err := db.AddItem(arg)
	if err != nil {
		return err
	}

	err = updateRecency(db, added)
	if err != nil {
		return err
	}

	return nil
}

// Attempts to call "go get" against an item with the highest score where the substr matches any items found in the db.
func goGetHighestScore(db *database.RummageDB, substr string, flags ...string) error {
	found, exists := db.EntryWithHighestScore(substr)
	if !exists {
		return errors.New("no entry was found that could match your argument")
	}

	err := attemptGoGet(found.Entry, flags...)
	if err != nil {
		return err
	}

	err = updateRecency(db, found)
	if err != nil {
		return err
	}

	return nil
}

// Tries to "go get" a package based on it's score in the database. If the argument passed has flashes, it will be treated as a regular go package
func Get(db *database.RummageDB, arg string, flags ...string) error {
	arg = strings.ToLower(arg)
	hasSlash, _ := utils.ParseForwardSlash(arg)
	// can safely assume if a "/" is parsed from the arg, it's more than likely an absolute pkg path
	if hasSlash {
		return goGetAddedItem(db, arg, flags...)
	}
	return goGetHighestScore(db, arg, flags...)
}
