package db

import (
	"fmt"
	"math"
	"time"
)

type RummageDBItem struct {
	Entry string
	// The "recency" score calculated and given to the item.
	Score float64
	// An int64 integer that represents the last time this value was accessed.
	// An "access" is considered anytime this value was last updated or the first time it was added.
	//
	// This time is in seconds, using the Unix epoch.
	LastAccessed int64
}

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

func (i *RummageDBItem) UpdateScore() *RummageDBItem {
	// TODO:
	// The score should meet these criteria
	// - be based on how long it's been since the user last accessed the db item,
	// - the score should be a float64 ranging between 0.01 (lowest match for search) - 0.99 (highest match for search)
	// - the longer it's been since last accessed, the lower the score should be (vice versa for recenctly accessed)

	return nil
}
