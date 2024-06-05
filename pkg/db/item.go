package db

import "time"

type RummageDBItem struct {
	Entry        string
	Score        int64
	LastAccessed time.Time
}
