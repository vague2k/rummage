package database_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/vague2k/rummage/pkg/database"
)

func TestAccess(t *testing.T) {
	t.Run("Initializing db does not error", func(t *testing.T) {
		tmp := t.TempDir()
		r, err := database.Init(tmp)
		if err != nil {
			t.Errorf("Could not open db: \n%s", err)
		}
		defer r.DB.Close()
	})
	t.Run("db returns correct dir and filepath", func(t *testing.T) {
		tmp := t.TempDir()
		got, err := database.Init(tmp)
		if err != nil {
			t.Errorf("Could not open db: \n%s", err)
		}
		defer got.DB.Close()

		expectedDir := filepath.Join(tmp, "rummage")
		expectedDBFile := filepath.Join(expectedDir, "rummage.db")

		switch true {
		case got.Dir != expectedDir:
			t.Errorf("Got %s, expected %s", got.Dir, expectedDir)
		case got.FilePath != expectedDBFile:
			t.Errorf("Got %s, expected %s", got.FilePath, expectedDBFile)
		}
	})
}

func TestAddItem(t *testing.T) {
	tmp := t.TempDir()
	r, err := database.Init(tmp)
	if err != nil {
		t.Errorf("Could not open db: \n%s", err)
	}
	defer r.DB.Close()

	t.Run("Can add item to db", func(t *testing.T) {
		_, err := r.AddItem("content")
		if err != nil {
			t.Errorf("Issue occured adding item to db: \n%s", err)
		}
	})

	t.Run("Assert added item has correct values", func(t *testing.T) {
		rows, err := r.DB.Query("SELECT * FROM items WHERE entry = 'content';")
		if err != nil {
			t.Errorf("Issue occured trying to select recently added item from db: \n%s", err)
		}

		var entry string
		var score float64
		var lastAccessed int64

		for rows.Next() {
			err := rows.Scan(&entry, &score, &lastAccessed)
			if err != nil {
				t.Errorf("Error occured when trying to scan rows: \n%s", err)
			}
		}

		switch true {
		case entry != "content":
			t.Errorf("Entry was %s, expected %s", entry, "content")
		case score != 1.0:
			t.Errorf("Score was %f, expected %f", score, 1.0)
		case lastAccessed != time.Now().Unix():
			t.Errorf("Created was %d, expected %v", lastAccessed, time.Now().Unix())
		}
	})
}

func TestSelectItem(t *testing.T) {
	tmp := t.TempDir()
	r, err := database.Init(tmp)
	if err != nil {
		t.Errorf("Could not open db: \n%s", err)
	}
	defer r.DB.Close()

	_, err = r.AddItem("firstitem")
	if err != nil {
		t.Errorf("Issue adding item to db: \n%s", err)
	}

	// checking the LastAccessed field is not neccessarily important for this test
	t.Run("Assert the added item exists", func(t *testing.T) {
		item, _ := r.SelectItem("firstitem")

		// checking the LastAccessed field is not neccessarily important for this check
		if item.Entry != "firstitem" {
			t.Errorf("Expected entry to be %s, but got %s.", "firstitem", item.Entry)
		}
		if item.Score != 1.0 {
			t.Errorf("Expected entry to be %f, but got %f.", 1.0, item.Score)
		}
	})

	t.Run("Assert added item is the only item in the db after attempting to add duplicate", func(t *testing.T) {
		_, err = r.AddItem("firstitem")
		if err != nil {
			t.Errorf("Issue adding item to db: \n%s", err)
		}

		rows, err := r.DB.Query("SELECT entry FROM items WHERE entry = ?", "firstitem")
		if err != nil {
			t.Errorf("Issue querying item from db: \n%s", err)
		}
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
	tmp := t.TempDir()
	r, err := database.Init(tmp)
	if err != nil {
		t.Errorf("Could not open db: \n%s", err)
	}

	originalItem, err := r.AddItem("firstitem")
	if err != nil {
		t.Errorf("Issue adding item to db: \n%s", err)
	}

	update := &database.RummageDBItem{
		Entry:        "updatedfirstitem",
		Score:        2.0,
		LastAccessed: time.Now().Unix(),
	}

	t.Run("Assert original item is actually updated", func(t *testing.T) {
		_, err := r.UpdateItem("firstitem", update)
		if err != nil {
			t.Errorf("Issue updating db item: \n%s", err)
		}

		_, err = r.DB.Query(`
            SELECT entry FROM items 
            WHERE score = ?`,
			2.0,
		)
		if err != nil {
			t.Errorf("Issue querying item from db: \n%s", err)
		}
	})

	t.Run("Returns pointer to updated item", func(t *testing.T) {
		updateItem, err := r.UpdateItem("firstitem", update)
		if err != nil {
			t.Errorf("Issue updating db item: \n%s", err)
		}

		switch true {
		case originalItem.Entry == updateItem.Entry:
			t.Errorf("Original entry %s and the updated entry are the same, expected %s", originalItem.Entry, update.Entry)
		case originalItem.Score == updateItem.Score:
			t.Errorf("Original score %f and the updated score are the same, expected %f", originalItem.Score, update.Score)
		}
	})
}
