package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/testutils"
)

func TestQuery(t *testing.T) {
	t.Run("Returns highest score entry", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/gorilla/mux", 2.0, 1)
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "mux")

		assert.Equal(t, "github.com/gorilla/mux \n", actual)
	})

	t.Run("Returns highest score entry and score", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/gorilla/mux", 2.0, 1)
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--score", "mux")

		assert.Equal(t, "github.com/gorilla/mux 2.00 \n", actual)
	})

	t.Run("Returns highest score entry, score and last-accessed", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/gorilla/mux", 2.0, 1725341448)
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--score", "--last-accessed", "mux")

		assert.Equal(t, "github.com/gorilla/mux 2.00 1725341448 \n", actual)
	})

	t.Run("Errors if no match was found", func(t *testing.T) {
		db := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
	})
}

func TestMultiQuery(t *testing.T) {
	t.Run("Returns multiple highest score entries", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/labstack/echo/v4", "github.com/user/echo")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/gorilla/mux", 2.0, 1)
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/labstack/echo/v4", 2.0, 1)
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "mux", "echo")

		assert.Equal(t, "github.com/gorilla/mux \ngithub.com/labstack/echo/v4 \n", actual)
	})

	t.Run("Returns multiple highest score entries and scores", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/labstack/echo/v4", "github.com/user/echo")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/gorilla/mux", 2.0, 1)
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/labstack/echo/v4", 2.0, 1)
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--score", "mux", "echo")

		assert.Equal(t, "github.com/gorilla/mux 2.00 \ngithub.com/labstack/echo/v4 2.00 \n", actual)
	})

	t.Run("Returns multiple highest score entries, scores and last-accessed fields", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/labstack/echo/v4", "github.com/user/echo")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/gorilla/mux", 2.0, 1725341448)
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/labstack/echo/v4", 2.0, 1725341448)
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--score", "--last-accessed", "mux", "echo")

		assert.Equal(t, "github.com/gorilla/mux 2.00 1725341448 \ngithub.com/labstack/echo/v4 2.00 1725341448 \n", actual)
	})

	t.Run("Errors if no multiple matches were found", func(t *testing.T) {
		db := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "mux", "echo")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement echo\n", actual)
	})
}

func TestExactQuery(t *testing.T) {
	t.Run("Returns exact entry", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "mux")

		assert.Equal(t, "github.com/gorilla/mux \n", actual)
	})

	t.Run("Returns exact entry and score", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--score", "mux")

		assert.Equal(t, "github.com/gorilla/mux 1.00 \n", actual)
	})

	t.Run("Returns exact entry, score and last accessed", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/gorilla/mux", 1.0, 1725341448)
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/user/mux", 1.0, 1725341448)
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--score", "--last-accessed", "mux")

		assert.Equal(t, "github.com/gorilla/mux 1.00 1725341448 \n", actual)
	})

	t.Run("Errors if no exact match was found", func(t *testing.T) {
		db := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--exact", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
	})
}

func TestMultiExactQuery(t *testing.T) {
	t.Run("Returns multiple exact entries", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/labstack/echo/v4", "github.com/user/echo")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "mux", "echo")

		assert.Equal(t, "github.com/gorilla/mux \ngithub.com/labstack/echo/v4 \n", actual)
	})

	t.Run("Returns multiple exact entries and scores", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/labstack/echo/v4", "github.com/user/echo")
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--exact", "--score", "mux", "echo")

		assert.Equal(t, "github.com/gorilla/mux 1.00 \ngithub.com/labstack/echo/v4 1.00 \n", actual)
	})

	t.Run("Returns multiple exact entries, scores and last-accessed", func(t *testing.T) {
		db := testutils.InMemDb(t)
		_, _, err := db.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/labstack/echo/v4", "github.com/user/echo")
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/gorilla/mux", 1.0, 1725341448)
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/user/mux", 1.0, 1725341448)
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/labstack/echo/v4", 1.0, 1725341448)
		assert.NoError(t, err)
		_, err = db.UpdateItem("github.com/user/echo", 1.0, 1725341448)
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--exact", "--score", "--last-accessed", "mux", "echo")

		assert.Equal(t, "github.com/gorilla/mux 1.00 1725341448 \ngithub.com/labstack/echo/v4 1.00 1725341448 \n", actual)
	})

	t.Run("Errors if no multiple exact matches were found", func(t *testing.T) {
		db := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "--exact", "mux", "echo")

		assert.Equal(t, "no match found with the given arguement mux\nno match found with the given arguement echo\n", actual)
	})
}
