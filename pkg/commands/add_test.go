package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	r := newDb(t)
	t.Run("Can add go packages", func(t *testing.T) {
		err := Add(r, "github.com/gorilla/mux", "github.com/charmbracelet/bubbletea")
		assert.NoError(t, err)
		items, err := r.ListItems()
		assert.NoError(t, err)
		assert.Len(t, items, 2)
	})
}
