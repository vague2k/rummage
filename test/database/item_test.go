package database_test

import (
	"testing"
	"time"

	"github.com/vague2k/rummage/pkg/database"
)

const (
	MINUTE = 60
	HOUR   = MINUTE * 60
	DAY    = HOUR * 24
	WEEK   = DAY * 7
)

var offset int64 // a small buffer of time for a more realistic-ish test

func TestRecalculateScore(t *testing.T) {
	t.Run("Calculate score within the hour", func(t *testing.T) {
		offset = MINUTE
		updated := &database.RummageDBItem{
			Entry:        "calculate",
			Score:        1.0,
			LastAccessed: time.Now().Unix() + offset,
		}

		expected := 4.0
		got := updated.RecalculateScore()
		if got != expected {
			t.Errorf("Expected %f, but got %f", expected, got)
		}
	})
	t.Run("Calculate score within the day", func(t *testing.T) {
		offset = HOUR
		updated := &database.RummageDBItem{
			Entry:        "calculate",
			Score:        1.0,
			LastAccessed: time.Now().Unix() + offset,
		}

		expected := 2.0
		got := updated.RecalculateScore()
		if got != expected {
			t.Errorf("Expected %f, but got %f", expected, got)
		}
	})
	t.Run("Calculate score within the week", func(t *testing.T) {
		offset = DAY
		updated := &database.RummageDBItem{
			Entry:        "calculate",
			Score:        1.0,
			LastAccessed: time.Now().Unix() + offset,
		}

		expected := 0.5
		got := updated.RecalculateScore()
		if got != expected {
			t.Errorf("Expected %f, but got %f", expected, got)
		}
	})
	t.Run("Calculate score past a week", func(t *testing.T) {
		offset = WEEK + WEEK
		updated := &database.RummageDBItem{
			Entry:        "calculate",
			Score:        1.0,
			LastAccessed: time.Now().Unix() + offset,
		}

		expected := 0.25
		got := updated.RecalculateScore()
		if got != expected {
			t.Errorf("Expected %f, but got %f", expected, got)
		}
	})
}
