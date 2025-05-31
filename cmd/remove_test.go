package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/testutils"
)

func TestRemove(t *testing.T) {
	t.Run("Properly removes item", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		_, err := db.AddItem(ctx, database.AddItemParams{
			Entry: "github.com/gorilla/mux",
		})
		assert.NoError(t, err)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "remove", "github.com/gorilla/mux")

		assert.Equal(t, "deleted github.com/gorilla/mux from the database\n", actual)
	})

	t.Run("Errors if item does not exist", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "remove", "doesnotexist")

		assert.Equal(t, "can't delete item with entry doesnotexist it does not exist\n", actual)
	})
}

func TestMultiRemove(t *testing.T) {
	t.Run("Properly removes item", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		pkgsToDelete := []string{"github.com/gorilla/mux", "github.com/user/mux"}
		for _, entry := range pkgsToDelete {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "remove", "github.com/gorilla/mux", "github.com/user/mux")

		assert.Equal(t, "deleted github.com/gorilla/mux from the database\ndeleted github.com/user/mux from the database\n", actual)
	})

	t.Run("Can delete items even if 1 item doesnt exist", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)
		pkgsToDelete := []string{"github.com/gorilla/mux", "github.com/user/mux", "doesnotexist"}
		for _, entry := range pkgsToDelete {
			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: entry,
			})
			assert.NoError(t, err)
		}

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "remove", "github.com/gorilla/mux", "doesnotexist", "github.com/user/mux")

		assert.Equal(t, "deleted github.com/gorilla/mux from the database\ncan't delete item with entry doesnotexist it does not exist\ndeleted github.com/user/mux from the database\n", actual)
	})

	t.Run("Errors if attempted items do not exist", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "remove", "doesnotexist", "doesnotexist2")

		assert.Equal(t, "can't delete item with entry doesnotexist it does not exist\ncan't delete item with entry doesnotexist2 it does not exist\n", actual)
	})
}
