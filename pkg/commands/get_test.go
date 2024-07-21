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
		assert.NoError(t, err)
	})

	t.Run("Errors when not a go package", func(t *testing.T) {
		test := "notagopackage"
		err := commands.AttemptGoGet(test)
		assert.Error(t, err)
	})

	t.Run("Can take u, t or x, args", func(t *testing.T) {
		tests := []struct {
			name  string
			flags []string
		}{
			{"takes -u", []string{"-u"}},
			{"takes -u -t", []string{"-u", "-t"}},
			{"takes -u -t -x", []string{"-u", "-t", "-x"}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				test := "github.com/gorilla/mux"
				err := commands.AttemptGoGet(test, tt.flags...)
				assert.NoError(t, err)
			})
		}
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
