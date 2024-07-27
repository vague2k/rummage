package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/database"
)

func newDb(t *testing.T) *database.RummageDB {
	r, err := database.Init(t.TempDir())
	assert.NoError(t, err)
	t.Cleanup(func() {
		r.DB.Close()
		r = nil
	})
	return r
}

func selectItem(r *database.RummageDB, s string) *database.RummageDBItem {
	found, _ := r.SelectItem(s)
	return found
}

func resetScore(t *testing.T, r *database.RummageDB, item *database.RummageDBItem) {
	reset := &database.RummageDBItem{
		Entry: item.Entry,
		Score: 1.0,
	}
	_, err := r.UpdateItem(item.Entry, reset)
	assert.NoError(t, err)
}
