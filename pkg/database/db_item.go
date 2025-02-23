package database

import (
	"time"
)

type RummageItem struct {
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
	_HOUR = 3600
	_DAY  = _HOUR * 24
	_WEEK = _DAY * 7
)

// Recalculates an item's score based on the last time it was last accessed.
func (i *RummageItem) RecalculateScore() float64 {
	var score float64

	duration := time.Now().Unix() - i.LastAccessed

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
