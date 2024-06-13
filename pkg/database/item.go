package database

import (
	"fmt"
	"time"
)

type RummageDBItem struct {
	// The name of the pkg to get
	Entry string
	// the score calculated based on recency. higher means more frequent uses.
	Score float64
	// An int64 integer that represents the last time this value was accessed.
	// An "access" is considered anytime this value was last updated or the first time it was added.
	//
	// This time is in seconds, using the Unix epoch.
	LastAccessed int64
}

const (
	HOUR = 3600
	DAY  = HOUR * 24
	WEEK = DAY * 7
)

// Returns the LastAccessed field as a string
func (i *RummageDBItem) LastAccessedAsString() string {
	s := fmt.Sprintf("%d", i.LastAccessed)
	return s
}

// Returns the Score field as a string
func (i *RummageDBItem) ScoreAsString() string {
	s := fmt.Sprintf("%f", i.Score)
	return s
}

// Recalculates an item's score based on the last time it was last accessed.
func (i *RummageDBItem) RecalculateScore() float64 {
	var score float64

	duration := i.LastAccessed - time.Now().Unix()

	// the older the time, the lower the score
	if duration < HOUR {
		score = i.Score * 4.0
	} else if duration < DAY {
		score = i.Score * 2.0
	} else if duration < WEEK {
		score = i.Score * 0.5
	} else {
		score = i.Score * 0.25
	}

	return score
}
