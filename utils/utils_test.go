package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/database"
)

// TODO: add test for ResemeblesGoPackage

const (
	MINUTE = 60
	HOUR   = MINUTE * 60
	DAY    = HOUR * 24
	WEEK   = DAY * 7
)

func TestRecalculateScore(t *testing.T) {
	tests := []struct {
		name          string
		offset        int64
		expectedScore float64
	}{
		{"Within the hour", MINUTE, 5.0},
		{"Within the day", HOUR, 3.0},
		{"Within the week", DAY, 0.5},
		{"Past a week", WEEK + WEEK, 0.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated := &database.RummageItem{
				Entry:        "calculate",
				Score:        1.0,
				Lastaccessed: time.Now().Unix() - tt.offset,
			}
			assert.Equal(t, tt.expectedScore, RecalculateScore(updated))
		})
	}
}
