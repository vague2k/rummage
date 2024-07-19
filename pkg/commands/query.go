package commands

import (
	"fmt"
	"strings"

	"github.com/vague2k/rummage/pkg/database"
)

func FindExactMatch(db *database.RummageDB, arg string) {
	arg = strings.ToLower(arg)
	found, exists := db.FindExactMatch(arg)
	if !exists {
		logger.Warn("No entry was found that could match your query.")
		return
	}
	info := printInfo(found)
	logger.Info("found exact match for ", arg, "\n", info)
}

func FindHighestScore(db *database.RummageDB, arg string) {
	arg = strings.ToLower(arg)
	found, exists := db.EntryWithHighestScore(arg)
	if !exists {
		logger.Warn("No entry was found that could match your query.")
		return
	}
	info := printInfo(found)
	logger.Info("found exact match for ", arg, "\n", info)
}

func printInfo(item *database.RummageDBItem) string {
	info := fmt.Sprintf("Entry: %s\nScore: %f\nLastAccessed: %d\n",
		item.Entry,
		item.Score,
		item.LastAccessed,
	)
	return info
}
