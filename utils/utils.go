package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/vague2k/rummage/pkg/database"
)

const (
	_HOUR = 3600
	_DAY  = _HOUR * 24
	_WEEK = _DAY * 7
)

// Gets the user's $GOPATH.
//
// Gets "user-home-dir/go" if $GOPATH does not exist.
func UserGoPath() string {
	if path := os.Getenv("GOPATH"); path != "" {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("\n", err)
	}

	return filepath.Join(home, "go")
}

func ResemblesGoPackage(entry string) error {
	regex := regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z0-9-!]+(/[a-zA-Z0-9-_\.!]+)+(/[vV]\d+)?$`)
	if !regex.MatchString(entry) {
		return fmt.Errorf("the item attempted to be added to the database does not resemble a go package")
	}
	return nil
}

// Recalculates an item's score based on the last time it was last accessed.
func RecalculateScore(item *database.RummageItem) float64 {
	var score float64

	duration := time.Now().Unix() - item.Lastaccessed

	// the older the time, the lower the score
	if duration < _HOUR {
		score = item.Score + 4.0
	} else if duration < _DAY {
		score = item.Score + 2.0
	} else if duration < _WEEK {
		score = item.Score * 0.5
	} else {
		score = item.Score * 0.25
	}

	return score
}
