package database_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		updated := &database.RummageItem{
			Entry:        "calculate",
			Score:        1.0,
			LastAccessed: time.Now().Unix() - offset,
		}

		assert.Equal(t, 5.0, updated.RecalculateScore())
	})
	t.Run("Calculate score within the day", func(t *testing.T) {
		offset = HOUR
		updated := &database.RummageItem{
			Entry:        "calculate",
			Score:        1.0,
			LastAccessed: time.Now().Unix() - offset,
		}

		assert.Equal(t, 3.0, updated.RecalculateScore())
	})
	t.Run("Calculate score within the week", func(t *testing.T) {
		offset = DAY
		updated := &database.RummageItem{
			Entry:        "calculate",
			Score:        1.0,
			LastAccessed: time.Now().Unix() - offset,
		}

		assert.Equal(t, 0.5, updated.RecalculateScore())
	})
	t.Run("Calculate score past a week", func(t *testing.T) {
		offset = WEEK + WEEK
		updated := &database.RummageItem{
			Entry:        "calculate",
			Score:        1.0,
			LastAccessed: time.Now().Unix() - offset,
		}

		assert.Equal(t, 0.25, updated.RecalculateScore())
	})
}
