package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/testutils"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Properly adds item",
			args:     []string{"add", "github.com/gorilla/mux"},
			expected: "added 1 package\n",
		},
		{
			name:     "Errors properly if attempted addition does not resemble go package",
			args:     []string{"add", "not-allowed1"},
			expected: "the item attempted to be added to the database does not resemble a go package\n",
		},
		{
			name:     "Properly adds multiple items",
			args:     []string{"add", "github.com/gorilla/mux", "github.com/labstack/echo/v4"},
			expected: "added 2 packages\n",
		},
		{
			name:     "Errors properly if all attempted additions do not resemble go packages",
			args:     []string{"add", "not-allowed1", "not-allowed2", "not-allowed3"},
			expected: "stopping the adding of items early...\nissue occured when adding item not-allowed1 to the db\nadded 0 packages\n",
		},
		{
			name:     "Errors properly midway through adding",
			args:     []string{"add", "github.com/gorilla/mux", "github.com/user/echo", "not-allowed1"},
			expected: "stopping the adding of items early...\nissue occured when adding item not-allowed1 to the db\nadded 2 packages\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := testutils.InMemDb(t)
			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tt.args...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
