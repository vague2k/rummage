package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/testutils"
)

func TestGetHighestScore(t *testing.T) {
	testutils.GoModTidy(t)
	t.Cleanup(func() {
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, err := db.AddItem("github.com/gorilla/mux")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db with higher score", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/charmbracelet/bubbletea", "github.com/charmbracelet/bubbles")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/charmbracelet/bubbletea", 5.0, time.Now().Unix())
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "bubble")

		assert.Contains(t, actual, "go: added github.com/charmbracelet/bubbletea")
		assert.NotContains(t, actual, "go: added github.com/charmbracelet/bubbles")
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db with -u flag", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, err := db.AddItem("github.com/gorilla/mux")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-u", "mux")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db with -t flag", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, err := db.AddItem("github.com/gorilla/mux")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-t", "mux")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		testutils.GoModTidy(t)
	})

	t.Run("Can get item from db with -x flag", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, err := db.AddItem("github.com/gorilla/mux")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-x", "mux")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		testutils.GoModTidy(t)
	})

	t.Run("Errors if item does not exist", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
		testutils.GoModTidy(t)
	})

	t.Run("Errors if item does not exist with -u flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-u", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
		testutils.GoModTidy(t)
	})

	t.Run("Errors if item does not exist with -t flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-t", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
		testutils.GoModTidy(t)
	})

	t.Run("Errors if item does not exist with -x flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

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
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/gofiber/fiber/v2")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "fiber")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		testutils.GoModTidy(t)
	})

	t.Run("Can get multiple item from db with higher score", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/charmbracelet/bubbletea", "github.com/charmbracelet/bubbles")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/charmbracelet/bubbletea", 5.0, time.Now().Unix())
		assert.NoError(t, err)
		_, _, err = db.AddMultiItems("golang.org/x/sync", "golang.org/x/net")
		assert.NoError(t, err)
		_, err = db.UpdateItem("golang.org/x/sync", 5.0, time.Now().Unix())
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
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/gofiber/fiber/v2")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-u", "mux", "fiber")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		testutils.GoModTidy(t)
	})

	t.Run("Can get multiple item from db with -t flag", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/gofiber/fiber/v2")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-t", "mux", "fiber")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		testutils.GoModTidy(t)
	})

	t.Run("Can get multiple item from db with -x flag", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/gofiber/fiber/v2")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-x", "mux", "fiber")

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		testutils.GoModTidy(t)
	})

	t.Run("Errors if multiple items do not exist", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "fiber")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement fiber\n", actual)
	})

	t.Run("Errors if multiple items do not exist with -u flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "-u", "fiber")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement fiber\n", actual)
	})

	t.Run("Errors if multiple items do not exist with -t flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "-t", "fiber")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement fiber\n", actual)
	})

	t.Run("Errors if multiple items do not exist with -x flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "mux", "-x", "fiber")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement fiber\n", actual)
	})
}

func TestGet(t *testing.T) {
	testutils.GoModTidy(t)
	t.Cleanup(func() {
		testutils.GoModTidy(t)
	})

	t.Run("Can get item", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "github.com/gorilla/mux")

		item, err := db.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Equal(t, 5.0, item.Score)
		assert.Equal(t, time.Now().Unix(), item.LastAccessed)
		testutils.GoModTidy(t)
	})

	t.Run("Can get item with -u flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-u", "github.com/gorilla/mux")

		item, err := db.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Equal(t, 5.0, item.Score)
		assert.Equal(t, time.Now().Unix(), item.LastAccessed)
		testutils.GoModTidy(t)
	})

	t.Run("Can get item with -t flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-t", "github.com/gorilla/mux")

		item, err := db.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Equal(t, 5.0, item.Score)
		assert.Equal(t, time.Now().Unix(), item.LastAccessed)
		testutils.GoModTidy(t)
	})

	t.Run("Can get item with -x flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-x", "github.com/gorilla/mux")

		item, err := db.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Equal(t, 5.0, item.Score)
		assert.Equal(t, time.Now().Unix(), item.LastAccessed)
		testutils.GoModTidy(t)
	})
}

func TestMultipleGet(t *testing.T) {
	testutils.GoModTidy(t)
	t.Cleanup(func() {
		testutils.GoModTidy(t)
	})

	t.Run("Can get items", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2")

		mux, err := db.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)
		fiber, err := db.SelectItem("github.com/gofiber/fiber/v2")
		assert.NoError(t, err)

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		assert.Equal(t, 5.0, mux.Score)
		assert.Equal(t, time.Now().Unix(), mux.LastAccessed)
		assert.Equal(t, 5.0, fiber.Score)
		assert.Equal(t, time.Now().Unix(), fiber.LastAccessed)
		testutils.GoModTidy(t)
	})

	t.Run("Can get items with -u flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-u", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2")

		mux, err := db.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)
		fiber, err := db.SelectItem("github.com/gofiber/fiber/v2")
		assert.NoError(t, err)

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		assert.Equal(t, 5.0, mux.Score)
		assert.Equal(t, time.Now().Unix(), mux.LastAccessed)
		assert.Equal(t, 5.0, fiber.Score)
		assert.Equal(t, time.Now().Unix(), fiber.LastAccessed)
		testutils.GoModTidy(t)
	})

	t.Run("Can get items with -t flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-t", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2")

		mux, err := db.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)
		fiber, err := db.SelectItem("github.com/gofiber/fiber/v2")
		assert.NoError(t, err)

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		assert.Equal(t, 5.0, mux.Score)
		assert.Equal(t, time.Now().Unix(), mux.LastAccessed)
		assert.Equal(t, 5.0, fiber.Score)
		assert.Equal(t, time.Now().Unix(), fiber.LastAccessed)
		testutils.GoModTidy(t)
	})

	t.Run("Can get items with -x flag", func(t *testing.T) {
		db := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "get", "-x", "github.com/gorilla/mux", "github.com/gofiber/fiber/v2")

		mux, err := db.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)
		fiber, err := db.SelectItem("github.com/gofiber/fiber/v2")
		assert.NoError(t, err)

		assert.Contains(t, actual, "go: added github.com/gorilla/mux")
		assert.Contains(t, actual, "go: added github.com/gofiber/fiber/v2")
		assert.Equal(t, 5.0, mux.Score)
		assert.Equal(t, time.Now().Unix(), mux.LastAccessed)
		assert.Equal(t, 5.0, fiber.Score)
		assert.Equal(t, time.Now().Unix(), fiber.LastAccessed)
		testutils.GoModTidy(t)
	})
}
