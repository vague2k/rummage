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
	testutils.GoModTidy(t)
	t.Cleanup(func() {
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		_, err := db.AddItem(ctx, database.AddItemParams{
			Entry: "github.com/gorilla/mux",
		})
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db with higher score", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		pkgs := []string{"github.com/charmbracelet/bubbletea", "github.com/charmbracelet/bubbles"}
		for _, entry := range pkgs {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}
		err := db.UpdateItem(ctx, database.UpdateItemParams{
			Entry:        "github.com/charmbracelet/bubbletea",
			Score:        5.0,
			Lastaccessed: time.Now().Unix(),
		})
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "bubble")

		assert.Contains(t, actual, "go: added github.com/charmbracelet/bubbletea")
		assert.NotContains(t, actual, "go: added github.com/charmbracelet/bubbles")
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db with -u flag", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		_, err := db.AddItem(ctx, database.AddItemParams{
			Entry: "github.com/gorilla/mux",
		})
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-u", "mux")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db with -t flag", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		_, err := db.AddItem(ctx, database.AddItemParams{
			Entry: "github.com/gorilla/mux",
		})
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-t", "mux")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db with -x flag", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		_, err := db.AddItem(ctx, database.AddItemParams{
			Entry: "github.com/gorilla/mux",
		})
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-x", "mux")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		testutils.GoModTidy(t)
	})

	t.Run("Errors if item does not exist", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
		testutils.GoModTidy(t)
	})

	t.Run("Errors if item does not exist with -u flag", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-u", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
		testutils.GoModTidy(t)
	})

	t.Run("Errors if item does not exist with -t flag", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-t", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
		testutils.GoModTidy(t)
	})

	t.Run("Errors if item does not exist with -x flag", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-x", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
		testutils.GoModTidy(t)
	})
}

func TestGetMultiHighestScore(t *testing.T) {
	testutils.GoModTidy(t)
	t.Cleanup(func() {
		testutils.GoModTidy(t)
	})

	t.Run("Can get multiple item from db", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		pkgs := []string{"github.com/gorilla/mux", "github.com/gofiber/fiber/v2"}
		for _, entry := range pkgs {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "fiber")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		testutils.GoModTidy(t)
	})

	t.Run("Can get multiple item from db with higher score", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		bubblePkgs := []string{"github.com/charmbracelet/bubbletea", "github.com/charmbracelet/bubbles"}
		for _, entry := range bubblePkgs {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}
		err := db.UpdateItem(ctx, database.UpdateItemParams{
			Entry:        "github.com/charmbracelet/bubbletea",
			Score:        5.0,
			Lastaccessed: time.Now().Unix(),
		})
		assert.NoError(t, err)

		goPkgs := []string{"golang.org/x/sync", "golang.org/x/net"}
		for _, entry := range goPkgs {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}
		err = db.UpdateItem(ctx, database.UpdateItemParams{
			Entry:        "golang.org/x/sync",
			Score:        5.0,
			Lastaccessed: time.Now().Unix(),
		})
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "bubble", "golang")

		assert.Contains(t, actual, "go: added github.com/charmbracelet/bubbletea")
		assert.Contains(t, actual, "go: added golang.org/x/sync")
		assert.NotContains(t, actual, "go: added github.com/charmbracelet/bubbles")
		assert.NotContains(t, actual, "go: added golang.org/x/net")
		testutils.GoModTidy(t)
	})

	t.Run("Can get multiple item from db with -u flag", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		bubblePkgs := []string{"github.com/gorilla/mux", "github.com/gofiber/fiber/v2"}
		for _, entry := range bubblePkgs {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-u", "mux", "fiber")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		testutils.GoModTidy(t)
	})

	t.Run("Can get multiple item from db with -t flag", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		bubblePkgs := []string{"github.com/gorilla/mux", "github.com/gofiber/fiber/v2"}
		for _, entry := range bubblePkgs {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-t", "mux", "fiber")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		testutils.GoModTidy(t)
	})

	t.Run("Can get multiple item from db with -x flag", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		bubblePkgs := []string{"github.com/gorilla/mux", "github.com/gofiber/fiber/v2"}
		for _, entry := range bubblePkgs {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-x", "mux", "fiber")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		testutils.GoModTidy(t)
	})

	t.Run("Errors if multiple items do not exist", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "fiber")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement fiber\n", actual)
	})

	t.Run("Errors if multiple items do not exist with -u flag", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "-u", "fiber")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement fiber\n", actual)
	})

	t.Run("Errors if multiple items do not exist with -t flag", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "-t", "fiber")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement fiber\n", actual)
	})

	t.Run("Errors if multiple items do not exist with -x flag", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "-x", "fiber")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement fiber\n", actual)
	})
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"Can get item", []string{"get", "github.com/gorilla/mux"}},
		{"Can get item with -u flag", []string{"get", "-u", "github.com/gorilla/mux"}},
		{"Can get item with -t flag", []string{"get", "-t", "github.com/gorilla/mux"}},
		{"Can get item with -x flag", []string{"get", "-x", "github.com/gorilla/mux"}},
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
		{"Can get items", []string{"get", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2"}},
		{"Can get items with -u flag", []string{"get", "-u", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2"}},
		{"Can get items with -t flag", []string{"get", "-t", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2"}},
		{"Can get items with -x flag", []string{"get", "-x", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2"}},
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
