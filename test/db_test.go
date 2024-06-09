package db_test

import (
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vague2k/rummage/pkg/db"
)

func TestAccess(t *testing.T) {
	t.Run("Accessing db does not error", func(t *testing.T) {
		tmp := t.TempDir()
		_, err := db.Access(tmp)
		if err != nil {
			t.Errorf("Could not open db: \n%s", err)
		}
	})
	t.Run("db returns correct dir and filepath", func(t *testing.T) {
		tmp := t.TempDir()
		got, err := db.Access(tmp)
		if err != nil {
			t.Errorf("Could not open db: \n%s", err)
		}

		expectedDir := filepath.Join(tmp, "rummage")
		expectedDBFile := filepath.Join(expectedDir, "db.rum")

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
	database, err := db.Access(tmp)
	if err != nil {
		t.Errorf("Could not open db: \n%s", err)
	}

	t.Run("Properly creates a db item", func(t *testing.T) {
		b := db.InternalCreateDBItem("content", 1.0, false)
		item := strings.Split(string(b), "\x00\x00")

		entry := item[0]
		score, _ := strconv.ParseFloat(item[1], 64)
		created, _ := strconv.ParseInt(item[2], 10, 0)

		switch true {
		case entry != "content":
			t.Errorf("Entry was %s, expected %s", entry, "content")
		case score != 1.0:
			t.Errorf("Score was %f, expected %f", score, 1.0)
		case created != time.Now().Unix():
			t.Errorf("Created was %d, expected %v", created, time.Now().Unix())
		}
	})

	t.Run("Properly adds item to db", func(t *testing.T) {
		item, err := database.AddItem("content")
		if err != nil {
			t.Errorf("Issue occured adding item to db: \n%s", err)
		}

		switch true {
		case item.Entry != "content":
			t.Errorf("Entry was %s, expected %s", item.Entry, "content")
		case item.Score != 1.0:
			t.Errorf("Score was %f, expected %f", item.Score, 1.0)
		case item.LastAccessed != time.Now().Unix():
			t.Errorf("Score was %d, expected %v", item.LastAccessed, time.Now().Unix())
		}
	})
}

func TestSelectItem(t *testing.T) {
	tmp := t.TempDir()
	db, err := db.Access(tmp)
	if err != nil {
		t.Errorf("Could not open db: \n%s", err)
	}

	_, err = db.AddItem("firstitem")
	if err != nil {
		t.Errorf("Issue adding item to db: \n%s", err)
	}

	items := db.ListItems()
	t.Run("Returns correct amount of items", func(t *testing.T) {
		if len(items) != 1 {
			t.Errorf("Expected ListItems to return %d items, instead got %d items.", 1, len(items))
		}
	})

	t.Run("Assert each item is what was added", func(t *testing.T) {
		// checking the LastAccessed field of each item is not neccessarily important for this check
		if items[0].Entry != "firstitem" {
			t.Errorf("Expected entry to be %s, but got %s.", "firstitem", items[0].Entry)
		}
		if items[0].Score != 1.0 {
			t.Errorf("Expected entry to be %f, but got %f.", 1.0, items[0].Score)
		}
	})
}

func TestUpdatedItem(t *testing.T) {
	tmp := t.TempDir()
	database, err := db.Access(tmp)
	if err != nil {
		t.Errorf("Could not open db: \n%s", err)
	}

	originalItem, err := database.AddItem("firstitem")
	if err != nil {
		t.Errorf("Issue adding item to db: \n%s", err)
	}

	t.Run("Returns correct updated item", func(t *testing.T) {
		update := &db.RummageDBItem{
			Entry:        "updatedfirstitem",
			Score:        2.0,
			LastAccessed: time.Now().Unix(),
		}

		updateItem, err := database.UpdateItem("firstitem", update)
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
