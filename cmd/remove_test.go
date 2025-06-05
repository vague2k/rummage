package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/testutils"
)

func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		addPkgs  []string
		args     []string
		expected string
	}{
		{
			name:     "Properly removes item",
			addPkgs:  []string{"github.com/gorilla/mux"},
			args:     []string{"remove", "github.com/gorilla/mux"},
			expected: "deleted github.com/gorilla/mux from the database\n",
		},
		{
			name:     "Errors if item does not exist",
			addPkgs:  nil,
			args:     []string{"remove", "doesnotexist"},
			expected: "can't delete item with entry doesnotexist it does not exist\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, ctx := testutils.InMemDb(t)
			for _, entry := range tc.addPkgs {
				_, err := db.AddItem(ctx, database.AddItemParams{Entry: entry})
				assert.NoError(t, err)
			}

			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tc.args...)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMultiRemove(t *testing.T) {
	tests := []struct {
		name     string
		addPkgs  []string
		args     []string
		expected string
	}{
		{
			name:    "Properly removes item",
			addPkgs: []string{"github.com/gorilla/mux", "github.com/user/mux"},
			args:    []string{"remove", "github.com/gorilla/mux", "github.com/user/mux"},
			expected: "deleted github.com/gorilla/mux from the database\n" +
				"deleted github.com/user/mux from the database\n",
		},
		{
			name:    "Can delete existing items even if 1 item doesn't exist",
			addPkgs: []string{"github.com/gorilla/mux", "github.com/user/mux"},
			args:    []string{"remove", "github.com/gorilla/mux", "doesnotexist", "github.com/user/mux"},
			expected: "deleted github.com/gorilla/mux from the database\n" +
				"can't delete item with entry doesnotexist it does not exist\n" +
				"deleted github.com/user/mux from the database\n",
		},
		{
			name:    "Errors if all items do not exist",
			addPkgs: nil,
			args:    []string{"remove", "doesnotexist", "doesnotexist2"},
			expected: "can't delete item with entry doesnotexist it does not exist\n" +
				"can't delete item with entry doesnotexist2 it does not exist\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, ctx := testutils.InMemDb(t)
			for _, entry := range tt.addPkgs {
				_, err := db.AddItem(ctx, database.AddItemParams{Entry: entry})
				assert.NoError(t, err)
			}

			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tt.args...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
