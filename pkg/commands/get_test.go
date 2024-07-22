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

	t.Run("Takes any permutation of '-u -t -x'", func(t *testing.T) {
		tests := []struct {
			name  string
			flags []string
		}{
			{"takes -u", []string{"-u"}},
			{"takes -t", []string{"-t"}},
			{"takes -x", []string{"-x"}},
			{"takes -u -t", []string{"-u", "-t"}},
			{"takes -u -x", []string{"-u", "-x"}},
			{"takes -t -u", []string{"-t", "-u"}},
			{"takes -t -x", []string{"-t", "-x"}},
			{"takes -x -u", []string{"-x", "-u"}},
			{"takes -x -t", []string{"-x", "-t"}},
			{"takes -u -t -x", []string{"-u", "-t", "-x"}},
			{"takes -u -x -t", []string{"-u", "-x", "-t"}},
			{"takes -t -u -x", []string{"-t", "-u", "-x"}},
			{"takes -t -x -u", []string{"-t", "-x", "-u"}},
			{"takes -x -u -t", []string{"-x", "-u", "-t"}},
			{"takes -x -t -u", []string{"-x", "-t", "-u"}},
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
