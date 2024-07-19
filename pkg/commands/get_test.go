package commands_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/commands"
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

func TestAttemptGoGet(t *testing.T) {
	t.Run("Successfully gets go package", func(t *testing.T) {
		test := "github.com/gorilla/mux"
		err := commands.AttemptGoGet(test)
		assert.Nil(t, err)
	})

	t.Run("Errors when not a go package", func(t *testing.T) {
		test := "notagopackage"
		err := commands.AttemptGoGet(test)
		assert.NotNil(t, err)
	})
}

func TestUpdateRecency(t *testing.T) {
	r := newDb(t)

	item, err := r.AddItem("item")
	if err != nil {
		t.Errorf("Issue adding item to db: \n%s", err)
	}
	t.Run("Properly increments score", func(t *testing.T) {
		got := commands.UpdateRecency(r, item)
		expected := 4.0

		assert.Equal(t, expected, got.Score)
	})
}
