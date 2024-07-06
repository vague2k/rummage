package database_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/testutils"
)

func TestAccess(t *testing.T) {
	t.Run("Initializing db does not error", func(t *testing.T) {
		r := testutils.DbInstance(t)
		defer r.DB.Close()
	})
	t.Run("db returns correct dir and filepath", func(t *testing.T) {
		tmp := t.TempDir()
		got, err := database.Init(tmp)
		testutils.CheckErr(t, err)
		defer got.DB.Close()

		expectedDir := filepath.Join(tmp, "rummage")
		expectedDBFile := filepath.Join(expectedDir, "rummage.db")

		testutils.AssertEquals(t, expectedDir, got.Dir)
		testutils.AssertEquals(t, expectedDBFile, got.FilePath)
	})
}

func TestAddItem(t *testing.T) {
	r := testutils.DbInstance(t)
	defer r.DB.Close()

	t.Run("Can add item to db", func(t *testing.T) {
		_, err := r.AddItem("content")
		testutils.CheckErr(t, err)
	})

	t.Run("Assert added item has correct values", func(t *testing.T) {
		rows, err := r.DB.Query("SELECT * FROM items WHERE entry = 'content';")
		testutils.CheckErr(t, err)

		var entry string
		var score float64
		var lastAccessed int64

		for rows.Next() {
			err := rows.Scan(&entry, &score, &lastAccessed)
			testutils.CheckErr(t, err)
		}

		testutils.AssertEquals(t, entry, "content")
		testutils.AssertEquals(t, score, 1.0)
		testutils.AssertEquals(t, lastAccessed, time.Now().Unix())
	})
}

func TestAddMultiItems(t *testing.T) {
	r := testutils.DbInstance(t)
	defer r.DB.Close()

	items, err := r.AddMultiItems("item1", "item2", "item3")
	testutils.CheckErr(t, err)

	t.Run("Returns expected amount of items", func(t *testing.T) {
		got := len(items)
		expected := 3

		testutils.AssertEquals(t, expected, got)
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
	r := testutils.DbInstance(t)
	defer r.DB.Close()

	_, err := r.AddItem("firstitem")
	testutils.CheckErr(t, err)

	// checking the LastAccessed field is not neccessarily important for this test
	t.Run("Assert the added item exists", func(t *testing.T) {
		item, _ := r.SelectItem("firstitem")

		// checking the LastAccessed field is not neccessarily important for this check
		testutils.AssertEquals(t, item.Entry, "firstitem")
		testutils.AssertEquals(t, item.Score, 1.0)
	})

	t.Run("Assert added item is the only item in the db after attempting to add duplicate", func(t *testing.T) {
		_, err = r.AddItem("firstitem")
		testutils.CheckErr(t, err)

		rows, err := r.DB.Query("SELECT entry FROM items WHERE entry = ?", "firstitem")
		testutils.CheckErr(t, err)
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
	r := testutils.DbInstance(t)
	defer r.DB.Close()

	originalItem, err := r.AddItem("firstitem")
	testutils.CheckErr(t, err)

	update := &database.RummageDBItem{
		Entry:        "updatedfirstitem",
		Score:        2.0,
		LastAccessed: time.Now().Unix(),
	}

	t.Run("Assert original item is actually updated", func(t *testing.T) {
		_, err := r.UpdateItem("firstitem", update)
		testutils.CheckErr(t, err)

		_, err = r.DB.Query(`
            SELECT entry FROM items 
            WHERE score = ?`,
			2.0,
		)
		testutils.CheckErr(t, err)
	})

	t.Run("Returns pointer to updated item", func(t *testing.T) {
		updateItem, err := r.UpdateItem("firstitem", update)
		testutils.CheckErr(t, err)

		testutils.AssertEquals(t, originalItem.Entry, updateItem.Entry)
		testutils.AssertNotEquals(t, originalItem.Score, updateItem.Score)
	})
}

func TestListItems(t *testing.T) {
	r := testutils.DbInstance(t)
	defer r.DB.Close()

	for i := range 5 {
		_, err := r.AddItem(fmt.Sprintf("item%d", i))
		testutils.CheckErr(t, err)
	}

	items, err := r.ListItems()
	testutils.CheckErr(t, err)

	t.Run("Returns expected amount of items", func(t *testing.T) {
		expected := 5
		got := len(items)

		testutils.AssertEquals(t, expected, got)
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
		r := testutils.DbInstance(t)
		defer r.DB.Close()

		for i := range 5 {
			name := fmt.Sprintf("item%d", i)
			_, err := r.AddItem(name)
			testutils.CheckErr(t, err)

			incrementedItemScore := database.RummageDBItem{
				Entry:        name,
				Score:        float64(i),
				LastAccessed: time.Now().Unix(),
			}

			_, err = r.UpdateItem(name, &incrementedItemScore)
			testutils.CheckErr(t, err)

		}
		got, _ := r.EntryWithHighestScore("it")
		expected := 4.0

		testutils.AssertEquals(t, expected, got.Score)
	})

	t.Run("Returns 1 item if it's the only item in the db", func(t *testing.T) {
		r := testutils.DbInstance(t)
		defer r.DB.Close()

		_, err := r.AddItem("item")
		testutils.CheckErr(t, err)

		got, _ := r.EntryWithHighestScore("it")
		expected := 1.0

		testutils.AssertEquals(t, expected, got.Score)
	})

	t.Run("Returns false if no match was found", func(t *testing.T) {
		r := testutils.DbInstance(t)
		defer r.DB.Close()

		got, gotExists := r.EntryWithHighestScore("it")
		expectedExists := false

		testutils.AssertNotNil(t, got)
		testutils.AssertEquals(t, expectedExists, gotExists)
	})
}

func TestDeleteItem(t *testing.T) {
	r := testutils.DbInstance(t)
	defer r.DB.Close()

	_, err := r.AddItem("item")
	testutils.CheckErr(t, err)

	t.Run("Assert item exists before deletion", func(t *testing.T) {
		got, _ := r.SelectItem("item")
		testutils.AssertNotNil(t, got)
	})

	t.Run("Assert item was deleted", func(t *testing.T) {
		_, err := r.DeleteItem("item")
		testutils.CheckErr(t, err)

		_, got := r.SelectItem("item")
		testutils.AssertFalse(t, got)
	})
}

func TestFindExactMatch(t *testing.T) {
	r := testutils.DbInstance(t)
	defer r.DB.Close()

	items, err := r.AddMultiItems("thisitemexistsfirst", "seconditem", "whathappenedhere")
	testutils.CheckErr(t, err)

	t.Run("Assert items exist before finding exact match", func(t *testing.T) {
		for _, item := range items {
			testutils.AssertNotNil(t, item)
		}
	})

	t.Run("Assert first item found regardless of score", func(t *testing.T) {
		found, _ := r.FindExactMatch("it")
		testutils.AssertEquals(t, "thisitemexistsfirst", found.Entry)
	})
}
