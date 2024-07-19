package database_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestAccess(t *testing.T) {
	t.Run("Initializing db does not error", func(t *testing.T) {
		newDb(t)
	})
	t.Run("db returns correct dir and filepath", func(t *testing.T) {
		tmp := t.TempDir()
		got, err := database.Init(tmp)
		assert.NoError(t, err)
		defer got.DB.Close()

		expectedDir := filepath.Join(tmp, "rummage")
		expectedDBFile := filepath.Join(expectedDir, "rummage.db")

		assert.Equal(t, expectedDir, got.Dir)
		assert.Equal(t, expectedDBFile, got.FilePath)
	})
}

func TestAddItem(t *testing.T) {
	r := newDb(t)

	t.Run("Can add item to db", func(t *testing.T) {
		_, err := r.AddItem("content")
		assert.NoError(t, err)
	})

	t.Run("Assert added item has correct values", func(t *testing.T) {
		rows, err := r.DB.Query("SELECT * FROM items WHERE entry = 'content';")
		assert.NoError(t, err)

		var entry string
		var score float64
		var lastAccessed int64

		for rows.Next() {
			err := rows.Scan(&entry, &score, &lastAccessed)
			assert.NoError(t, err)
		}

		assert.Equal(t, entry, "content")
		assert.Equal(t, score, 1.0)
		assert.Equal(t, lastAccessed, time.Now().Unix())
	})
}

func TestAddMultiItems(t *testing.T) {
	r := newDb(t)

	items, err := r.AddMultiItems("item1", "item2", "item3")
	assert.NoError(t, err)

	t.Run("Returns expected amount of items", func(t *testing.T) {
		got := len(items)
		expected := 3

		assert.Equal(t, expected, got)
	})

	t.Run("Assert each item is correctly typed", func(t *testing.T) {
		for _, item := range items {
			var check interface{} = item
			if value, ok := check.(database.RummageDBItem); !ok {
				t.Errorf("The item %v is not of type database.RummageDBItem", value)
			}
		}
	})
}

func TestSelectItem(t *testing.T) {
	r := newDb(t)

	_, err := r.AddItem("firstitem")
	assert.NoError(t, err)

	// checking the LastAccessed field is not neccessarily important for this test
	t.Run("Assert the added item exists", func(t *testing.T) {
		item, _ := r.SelectItem("firstitem")

		// checking the LastAccessed field is not neccessarily important for this check
		assert.Equal(t, item.Entry, "firstitem")
		assert.Equal(t, item.Score, 1.0)
	})

	t.Run("Assert added item is the only item in the db after attempting to add duplicate", func(t *testing.T) {
		_, err = r.AddItem("firstitem")
		assert.NoError(t, err)

		rows, err := r.DB.Query("SELECT entry FROM items WHERE entry = ?", "firstitem")
		assert.NoError(t, err)
		defer rows.Close()

		count := 0
		for rows.Next() {
			count++
		}

		if count != 1 {
			t.Errorf("Expected 1 item in the db, but found %d.", count)
		}
	})

	t.Run("Assert an item does not exist", func(t *testing.T) {
		entryShouldNotExist := "somethingthatshouldntexist"
		_, exists := r.SelectItem(entryShouldNotExist)

		if exists {
			t.Errorf("Expected item with entry '%s' to not exist, but it actually does.", entryShouldNotExist)
		}
	})
}

func TestUpdatedItem(t *testing.T) {
	r := newDb(t)

	originalItem, err := r.AddItem("firstitem")
	assert.NoError(t, err)

	update := &database.RummageDBItem{
		Entry:        "updatedfirstitem",
		Score:        2.0,
		LastAccessed: time.Now().Unix(),
	}

	t.Run("Assert original item is actually updated", func(t *testing.T) {
		_, err := r.UpdateItem("firstitem", update)
		assert.NoError(t, err)

		_, err = r.DB.Query(`
            SELECT entry FROM items 
            WHERE score = ?`,
			2.0,
		)
		assert.NoError(t, err)
	})

	t.Run("Returns pointer to updated item", func(t *testing.T) {
		updateItem, err := r.UpdateItem("firstitem", update)
		assert.NoError(t, err)

		assert.Equal(t, originalItem.Entry, updateItem.Entry)
		assert.NotEqual(t, originalItem.Score, updateItem.Score)
	})
}

func TestListItems(t *testing.T) {
	r := newDb(t)

	for i := range 5 {
		_, err := r.AddItem(fmt.Sprintf("item%d", i))
		assert.NoError(t, err)
	}

	items, err := r.ListItems()
	assert.NoError(t, err)

	t.Run("Returns expected amount of items", func(t *testing.T) {
		expected := 5
		got := len(items)

		assert.Equal(t, expected, got)
	})

	t.Run("Assert each item is correctly typed", func(t *testing.T) {
		for _, item := range items {
			var check interface{} = item
			if value, ok := check.(database.RummageDBItem); !ok {
				t.Errorf("The item %v is not of type database.RummageDBItem", value)
			}
		}
	})
}

func TestEntryWithHighestScore(t *testing.T) {
	t.Run("Returns highest score if multiple entries are found", func(t *testing.T) {
		r := newDb(t)

		for i := range 5 {
			name := fmt.Sprintf("item%d", i)
			_, err := r.AddItem(name)
			assert.NoError(t, err)

			incrementedItemScore := database.RummageDBItem{
				Entry:        name,
				Score:        float64(i),
				LastAccessed: time.Now().Unix(),
			}

			_, err = r.UpdateItem(name, &incrementedItemScore)
			assert.NoError(t, err)

		}
		got, _ := r.EntryWithHighestScore("it")
		expected := 4.0

		assert.Equal(t, expected, got.Score)
	})

	t.Run("Returns 1 item if it's the only item in the db", func(t *testing.T) {
		r := newDb(t)

		_, err := r.AddItem("item")
		assert.NoError(t, err)

		got, _ := r.EntryWithHighestScore("it")
		expected := 1.0

		assert.Equal(t, expected, got.Score)
	})

	t.Run("Returns false if no match was found", func(t *testing.T) {
		r := newDb(t)

		got, gotExists := r.EntryWithHighestScore("it")
		expectedExists := false

		assert.Nil(t, got)
		assert.Equal(t, expectedExists, gotExists)
	})
}

func TestDeleteItem(t *testing.T) {
	r := newDb(t)

	_, err := r.AddItem("item")
	assert.NoError(t, err)

	t.Run("Assert item exists before deletion", func(t *testing.T) {
		got, _ := r.SelectItem("item")
		assert.NotNil(t, got)
	})

	t.Run("Assert item was deleted", func(t *testing.T) {
		_, err := r.DeleteItem("item")
		assert.NoError(t, err)

		_, got := r.SelectItem("item")
		assert.False(t, got)
	})
}

func TestFindExactMatch(t *testing.T) {
	r := newDb(t)

	items, err := r.AddMultiItems("thisitemexistsfirst", "seconditem", "whathappenedhere")
	assert.NoError(t, err)

	t.Run("Assert items exist before finding exact match", func(t *testing.T) {
		for _, item := range items {
			assert.NotNil(t, item)
		}
	})

	t.Run("Assert first item found regardless of score", func(t *testing.T) {
		found, _ := r.FindExactMatch("it")
		assert.Equal(t, "thisitemexistsfirst", found.Entry)
	})
}
