package commands

import (
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/ui"
)

// Prompts to user to confirm wether or not they want to delete all items from the database.
func PromptDeleteAll(db *database.RummageDB) {
	choice := ui.YesNoPrompt("Are you sure you want to delete all items from the database?", false)
	if !choice {
		return
	}
	err := db.DeleteAllItems()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Warn("Deleted ALL from the database...")
}
