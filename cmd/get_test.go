package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/testutils"
)

// NOTE: To fix the issue where github workflow has a discrepency where
// LastAccessed field is asserted to be time.Now().Unix() fails because the value is off by 1,
// We only can only check the first 8 digits of the Unix Epoch intead of all 10.
// these tests do not need such preciseness
// See testutils.First8()

func TestGetHighestScore(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Get item with highest score",
			args: []string{"get", "bubble"},
		},
		{
			name: "Get item with highest score with -u",
			args: []string{"get", "-u", "bubble"},
		},
		{
			name: "Get item with highest score with -t",
			args: []string{"get", "-t", "bubble"},
		},
		{
			name: "Get item with highest score with -x",
			args: []string{"get", "-x", "bubble"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.GoModTidy(t)
			t.Cleanup(func() {
				testutils.GoModTidy(t)
			})

			pkgs := []string{"github.com/charmbracelet/bubbletea", "github.com/charmbracelet/bubbles"}
			db, ctx := testutils.InMemDb(t)

			for _, entry := range pkgs {
				_, err := db.AddItem(ctx, database.AddItemParams{
					Entry: entry,
				})
				assert.NoError(t, err)
				err = db.UpdateItem(ctx, database.UpdateItemParams{
					Entry: "github.com/charmbracelet/bubbletea",
					Score: 5.0,
				})
				assert.NoError(t, err)
			}

			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tt.args...)

			assert.Contains(t, actual, "go: added github.com/charmbracelet/bubbletea")
			assert.NotContains(t, actual, "go: added github.com/charmbracelet/bubbles")
			testutils.GoModTidy(t)
		})
	}
}

func TestGetHighestScoreErrors(t *testing.T) {
	testutils.GoModTidy(t)
	t.Cleanup(func() {
		testutils.GoModTidy(t)
	})

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Errors if item does not exist",
			args:     []string{"get", "mux"},
			expected: "no match found with the given arguement mux\n",
		},
		{
			name:     "Errors if item does not exist with -u",
			args:     []string{"get", "-u", "mux"},
			expected: "no match found with the given arguement mux\n",
		},
		{
			name:     "Errors if item does not exist with -t",
			args:     []string{"get", "-t", "mux"},
			expected: "no match found with the given arguement mux\n",
		},
		{
			name:     "Errors if item does not exist with -x",
			args:     []string{"get", "-x", "mux"},
			expected: "no match found with the given arguement mux\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := testutils.InMemDb(t)

			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tt.args...)

			assert.Equal(t, "no match found with the given arguement mux\n", actual)
			testutils.GoModTidy(t)
		})
	}
}

func TestGetMultiHighestScore(t *testing.T) {
	testutils.GoModTidy(t)
	t.Cleanup(func() {
		testutils.GoModTidy(t)
	})

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Get multiple items with highest score",
			args: []string{"get", "bubble", "mux"},
		},
		{
			name: "Get multiple items with highest score with -u",
			args: []string{"get", "-u", "bubble", "mux"},
		},
		{
			name: "Get multiple items with highest score with -t",
			args: []string{"get", "-t", "bubble", "mux"},
		},
		{
			name: "Get multiple items with highest score with -x",
			args: []string{"get", "-x", "bubble", "mux"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.GoModTidy(t)
			t.Cleanup(func() {
				testutils.GoModTidy(t)
			})
			db, ctx := testutils.InMemDb(t)
			pkgs := []string{"github.com/charmbracelet/bubbletea", "github.com/charmbracelet/bubbles", "github.com/gorilla/mux"}
			for _, entry := range pkgs {
				_, err := db.AddItem(ctx, database.AddItemParams{
					Entry: entry,
				})
				assert.NoError(t, err)
				err = db.UpdateItem(ctx, database.UpdateItemParams{
					Entry: "github.com/charmbracelet/bubbletea",
					Score: 5.0,
				})
				assert.NoError(t, err)
			}

			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tt.args...)

			assert.Contains(t, actual, "go: added github.com/charmbracelet/bubbletea")
			assert.Contains(t, actual, "go: added github.com/gorilla/mux")
			assert.NotContains(t, actual, "go: added github.com/charmbracelet/bubbles")
			testutils.GoModTidy(t)
		})
	}
}

func TestGetMultiHighestScoreErrors(t *testing.T) {
	testutils.GoModTidy(t)
	t.Cleanup(func() {
		testutils.GoModTidy(t)
	})

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Errors if item does not exist",
			args: []string{"get", "bubbles", "mux"},
		},
		{
			name: "Errors if item does not exist with -u",
			args: []string{"get", "-u", "bubbles", "mux"},
		},
		{
			name: "Errors if item does not exist with -t",
			args: []string{"get", "-t", "bubbles", "mux"},
		},
		{
			name: "Errors if item does not exist with -x",
			args: []string{"get", "-x", "bubbles", "mux"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := testutils.InMemDb(t)

			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tt.args...)

			assert.Equal(t, "no match found with the given arguement bubbles\nno match found with the given arguement mux\n", actual)
			testutils.GoModTidy(t)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Can get item",
			args: []string{"get", "github.com/gorilla/mux"},
		},
		{
			name: "Can get item with -u flag",
			args: []string{"get", "-u", "github.com/gorilla/mux"},
		},
		{
			name: "Can get item with -t flag",
			args: []string{"get", "-t", "github.com/gorilla/mux"},
		},
		{
			name: "Can get item with -x flag",
			args: []string{"get", "-x", "github.com/gorilla/mux"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testutils.GoModTidy(t)
			t.Cleanup(func() {
				testutils.GoModTidy(t)
			})
			db, ctx := testutils.InMemDb(t)

			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tc.args...)

			item, err := db.SelectItem(ctx, "github.com/gorilla/mux")
			assert.NoError(t, err)

			assert.Contains(t, actual, "go: added github.com/gorilla/mux")
			assert.Equal(t, 5.0, item.Score)
			assert.Equal(t, testutils.First8(time.Now().Unix()), testutils.First8(item.Lastaccessed))
			testutils.GoModTidy(t)
		})
	}
}

func TestMultipleGet(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Can get items",
			args: []string{"get", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2"},
		},
		{
			name: "Can get items with -u flag",
			args: []string{"get", "-u", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2"},
		},
		{
			name: "Can get items with -t flag",
			args: []string{"get", "-t", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2"},
		},
		{
			name: "Can get items with -x flag",
			args: []string{"get", "-x", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testutils.GoModTidy(t)
			t.Cleanup(func() {
				testutils.GoModTidy(t)
			})
			db, ctx := testutils.InMemDb(t)

			cmd := NewRootCmd(db)
			actual := testutils.Execute(cmd, tc.args...)

			pkgs := []string{
				"github.com/gorilla/mux",
				"github.com/gofiber/fiber/v2",
			}

			for _, pkg := range pkgs {
				item, err := db.SelectItem(ctx, pkg)
				assert.NoError(t, err)
				assert.Contains(t, actual, fmt.Sprintf("go: added %s", pkg))
				assert.Equal(t, 5.0, item.Score)
				assert.Equal(t, testutils.First8(time.Now().Unix()), testutils.First8(item.Lastaccessed))
			}

			testutils.GoModTidy(t)
		})
	}
}
