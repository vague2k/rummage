package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vague2k/rummage/utils"
)

//go:embed schema.sql
var schema string

const (
	_HOUR = 3600
	_DAY  = _HOUR * 24
	_WEEK = _DAY * 7
)

// Recalculates an item's score based on the last time it was last accessed.
func (i *RummageItem) RecalculateScore() float64 {
	var score float64

	duration := time.Now().Unix() - i.Lastaccessed

	// the older the time, the lower the score
	if duration < _HOUR {
		score = i.Score + 4.0
	} else if duration < _DAY {
		score = i.Score + 2.0
	} else if duration < _WEEK {
		score = i.Score * 0.5
	} else {
		score = i.Score * 0.25
	}

	return score
}

func Init(path string) (*Queries, error) {
	if path == "" {
		dataDir := utils.UserDataDir()
		path = dataDir
	}

	var dir string
	var dbFile string
	if path == ":memory:" {
		dbFile = ":memory:"
	} else {
		dir = filepath.Join(path, "rummage")
		dbFile = filepath.Join(dir, "rummage.db")
		err := os.MkdirAll(dir, 0o777)
		if err != nil {
			return nil, fmt.Errorf("could not create db dir: \n%s", err)
		}
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not init rummage db: \n%s", err)
	}

	if _, err := db.ExecContext(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("could not create 'items' table in rummage db: \n%s", err)
	}

	return New(db), nil
}
