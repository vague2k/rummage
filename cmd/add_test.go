package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/testutils"
)

func TestAdd(t *testing.T) {
	t.Run("Properly adds item", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "add", "github.com/gorilla/mux")

		assert.Equal(t, "added 1 package\n", actual)
	})
	t.Run("Errors properly if attempted addition does not resemble go package", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "add", "not-allowed1")

		assert.Equal(t, "the item attempted to be added to the database does not resemble a go package\n", actual)
	})
}

func TestMultipleAdd(t *testing.T) {
	t.Run("Properly adds items", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "add", "github.com/gorilla/mux", "github.com/labstack/echo/v4")

		assert.Equal(t, "added 2 packages\n", actual)
	})

	t.Run("Errors properly if first attempted additions do not resemble go package", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "add", "not-allowed1", "not-allowed2", "not-allowed3")

		assert.Equal(t, "stopping the adding of items early...\nissue occured when adding item not-allowed1 to the db\nadded 0 packages\n", actual)
	})

	t.Run("Errors properly midway through adding", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "add", "github.com/gorilla/mux", "github.com/user/echo", "not-allowed1")

		assert.Equal(t, "stopping the adding of items early...\nissue occured when adding item not-allowed1 to the db\nadded 2 packages\n", actual)
	})
}
