package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	// test setup
	r := newDb(t)
	if _, err := r.AddMultiItems("github.com/gorilla/mux"); err != nil {
		t.Error(err)
	}

	// actual test cases
	t.Run("Can get package in database without flags", func(t *testing.T) {
		err := Get(r, "mux")
		item := selectItem(r, "github.com/gorilla/mux")
		assert.NoError(t, err)
		assert.Equal(t, 4.0, item.Score)
		resetScore(t, r, item)
	})

	t.Run("Can get package in database with flags", func(t *testing.T) {
		err := Get(r, "mux", "-u")
		item := selectItem(r, "github.com/gorilla/mux")
		assert.NoError(t, err)
		assert.Equal(t, 4.0, item.Score)
		resetScore(t, r, item)
	})

	t.Run("Can get package not in database without flags", func(t *testing.T) {
		err := Get(r, "github.com/charmbracelet/bubbles")
		item := selectItem(r, "github.com/charmbracelet/bubbles")
		assert.NoError(t, err)
		assert.Equal(t, 4.0, item.Score)
		resetScore(t, r, item)
	})

	t.Run("Can get package not in database with flags", func(t *testing.T) {
		err := Get(r, "github.com/charmbracelet/bubbles", "-u")
		item := selectItem(r, "github.com/charmbracelet/bubbles")
		assert.NoError(t, err)
		assert.Equal(t, 4.0, item.Score)
		resetScore(t, r, item)
	})

	t.Run("Errors when no entry could be found", func(t *testing.T) {
		err := Get(r, "shouldnotexist")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no entry was found that could match your argument")
	})

	t.Run("Errors when package can't be got... without flags", func(t *testing.T) {
		if _, err := r.AddItem("shouldnotexist"); err != nil {
			t.Error(err)
		}
		err := Get(r, "shouldnotexist")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "go: malformed module path")
	})

	t.Run("Errors when package can't be got... with flags", func(t *testing.T) {
		err := Get(r, "shouldnotexist", "-u")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "go: malformed module path")
	})
}
