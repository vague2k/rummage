package database_test

import (
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/database"
)

// Spin up an in memory db (since we're using sqlite3) for quick testing
func inMemDb(t *testing.T) *database.RummageDb {
	r, err := database.Init(":memory:")
	assert.NoError(t, err)
	t.Cleanup(func() {
		r.Close()
		r = nil
	})
	return r
}

// Actual test cases should are in seperate t.Run() instances, unless the test is reasonable concise enough to
// be put under it's own function.
//
// For the most part, if you see function calls before a t.Run(), it's more than likely a setup for those test cases.

func TestInit(t *testing.T) {
	tmp := t.TempDir()
	r, err := database.Init(tmp)
	assert.NoError(t, err)

	expectedDir := filepath.Join(tmp, "rummage")
	expectedDBFile := filepath.Join(expectedDir, "rummage.db")

	assert.NotNil(t, r)
	assert.NotEmpty(t, r.Dir)
	assert.NotEmpty(t, r.FilePath)
	assert.Equal(t, expectedDir, r.Dir)
	assert.Equal(t, expectedDBFile, r.FilePath)
}

func TestAddItem(t *testing.T) {
	r := inMemDb(t)

	t.Run("Can add item that resembles a go package", func(t *testing.T) {
		item, err := r.AddItem("github.com/gorilla/mux")
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "github.com/gorilla/mux", item.Entry)
		assert.Equal(t, 1.0, item.Score)
		assert.Equal(t, time.Now().Unix(), item.LastAccessed)
	})

	// we can check if it's the same item by comparing the item.LastAccessed field offset by 1 (second) in both directions
	t.Run("Returns existing item", func(t *testing.T) {
		item, err := r.AddItem("github.com/gorilla/mux")
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "github.com/gorilla/mux", item.Entry)
		assert.Equal(t, 1.0, item.Score)
		assert.NotEqual(t, time.Now().Unix(), item.LastAccessed+1)
		assert.NotEqual(t, time.Now().Unix(), item.LastAccessed-1)
	})

	t.Run("Errors if item doesn't resemble a go package", func(t *testing.T) {
		item, err := r.AddItem("notagopackage")
		assert.ErrorContains(t, err, "the item attempted to be added to the database does not resemble a go package")
		assert.Nil(t, item)
	})
}

func TestAddMultiItems(t *testing.T) {
	t.Run("Can add item(s) that resembles a go package", func(t *testing.T) {
		r := inMemDb(t)
		items, amtAdded, err := r.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/doesntexist/mux")
		assert.NoError(t, err)
		assert.Equal(t, 3, amtAdded)
		for _, item := range items {
			assert.Regexp(t, regexp.MustCompile(`^github\.com/(gorilla|user|doesntexist)/mux$`), item.Entry)
			assert.Equal(t, 1.0, item.Score)
			assert.Equal(t, time.Now().Unix(), item.LastAccessed)
		}
	})

	t.Run("Returns 0 added if items already exist in db. 0/3", func(t *testing.T) {
		r := inMemDb(t)
		_, amtAdded, err := r.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/doesntexist/mux")
		assert.NoError(t, err)
		assert.Equal(t, 3, amtAdded)

		_, amtAgain, err := r.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/doesntexist/mux")
		assert.NoError(t, err)
		assert.Equal(t, 0, amtAgain)
	})

	t.Run("Errors if item doesn't resemble a go package", func(t *testing.T) {
		r := inMemDb(t)
		items, amtAdded, err := r.AddMultiItems("notagopackage", "notanothergopackage")
		assert.ErrorContains(t, err, "issue occured when adding item notagopackage to the db")
		assert.Nil(t, items)
		assert.Equal(t, 0, amtAdded)
	})

	t.Run("Returns correct amount of added packages. 2/3", func(t *testing.T) {
		r := inMemDb(t)
		items, amtAdded, err := r.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "notagopackage")
		assert.Error(t, err)
		assert.NotNil(t, items)
		assert.Equal(t, 2, amtAdded)
	})
}

func TestSelectItem(t *testing.T) {
	r := inMemDb(t)
	_, err := r.AddItem("github.com/gorilla/mux")
	assert.NoError(t, err)

	t.Run("Can select existing item", func(t *testing.T) {
		item, err := r.SelectItem("github.com/gorilla/mux")
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "github.com/gorilla/mux", item.Entry)
		assert.Equal(t, 1.0, item.Score)
		assert.Equal(t, time.Now().Unix(), item.LastAccessed)
	})

	t.Run("Errors if item does not exist", func(t *testing.T) {
		item, err := r.SelectItem("doesnotexist")
		assert.ErrorContains(t, err, "the item with entry doesnotexist does not exist")
		assert.Nil(t, item)
	})
}

func TestUpdateItem(t *testing.T) {
	r := inMemDb(t)
	_, err := r.AddItem("github.com/gorilla/mux")
	assert.NoError(t, err)

	t.Run("Can update existing item", func(t *testing.T) {
		item, err := r.UpdateItem("github.com/gorilla/mux", 4.0, 123)
		assert.NoError(t, err)
		assert.Equal(t, "github.com/gorilla/mux", item.Entry)
		assert.Equal(t, 4.0, item.Score)
		assert.Equal(t, int64(123), item.LastAccessed)
	})

	t.Run("Errors if item does not exist", func(t *testing.T) {
		item, err := r.UpdateItem("shouldnotexist", 0.0, 1)
		assert.ErrorContains(t, err, "the item with entry shouldnotexist is attempted to be updated but does not exist")
		assert.Nil(t, item)
	})
}

func TestListItems(t *testing.T) {
	r := inMemDb(t)
	_, _, err := r.AddMultiItems("github.com/gorilla/mux", "github.com/gofiber/fiber/v2")
	assert.NoError(t, err)

	items, err := r.ListItems()
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestDeleteItem(t *testing.T) {
	r := inMemDb(t)
	_, err := r.AddItem("github.com/gorilla/mux")
	assert.NoError(t, err)

	t.Run("Can delete existing item", func(t *testing.T) {
		item, err := r.DeleteItem("github.com/gorilla/mux")
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "github.com/gorilla/mux", item.Entry)
		assert.Equal(t, 1.0, item.Score)
		assert.Equal(t, time.Now().Unix(), item.LastAccessed)

		item, err = r.SelectItem("github.com/gorilla/mux")
		assert.Error(t, err)
		assert.Nil(t, item)
	})

	t.Run("Errors if item does not exist", func(t *testing.T) {
		item, err := r.DeleteItem("doesnotexist")
		assert.ErrorContains(t, err, "can't delete item with entry doesnotexist it does not exist")
		assert.Nil(t, item)
	})
}

func TestEntryWithHighestScore(t *testing.T) {
	r := inMemDb(t)
	_, _, err := r.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/doesntexist/mux")
	assert.NoError(t, err)
	_, err = r.UpdateItem("github.com/gorilla/mux", 2.0, 2)
	assert.NoError(t, err)
	_, err = r.UpdateItem("github.com/user/mux", 10.0, 10)
	assert.NoError(t, err)
	_, err = r.UpdateItem("github.com/doesntexist/mux", 5.0, 5)
	assert.NoError(t, err)

	t.Run("Can get entry with highest score", func(t *testing.T) {
		item, err := r.EntryWithHighestScore("mux")
		assert.NoError(t, err)
		assert.NotNil(t, item)

		assert.Equal(t, "github.com/user/mux", item.Entry)
		assert.Equal(t, 10.0, item.Score)
		assert.Equal(t, int64(10), item.LastAccessed)
	})

	t.Run("Errors if no match was found", func(t *testing.T) {
		item, err := r.EntryWithHighestScore("doesnotexist")
		assert.ErrorContains(t, err, "no match found with the given arguement doesnotexist")
		assert.Nil(t, item)
	})
}

func TestFindExactMatch(t *testing.T) {
	r := inMemDb(t)
	_, _, err := r.AddMultiItems("github.com/gorilla/mux", "github.com/user/mux", "github.com/doesntexist/mux")
	assert.NoError(t, err)

	t.Run("Can find first occurence of substr (exact match)", func(t *testing.T) {
		item, err := r.FindExactMatch("mux")
		assert.NoError(t, err)
		assert.NotNil(t, item)

		assert.Equal(t, "github.com/gorilla/mux", item.Entry)
		assert.Equal(t, 1.0, item.Score)
		assert.Equal(t, time.Now().Unix(), item.LastAccessed)
	})

	t.Run("Errors if no exact match was found", func(t *testing.T) {
		item, err := r.FindExactMatch("doesnotexist")
		assert.ErrorContains(t, err, "no match found with the given arguement doesnotexist")
		assert.Nil(t, item)
	})
}
