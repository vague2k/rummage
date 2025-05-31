package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/testutils"
)

func TestPopulate(t *testing.T) {
	t.Run("Can populate db with 3 out of 3 packages", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "populate", "--dir="+testutils.Mock3outof3pkgs(t))

		assert.Equal(t, "added 3 packages\n", actual)
	})

	t.Run("Can populate db with 1 out of 3 packages", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "populate", "--dir="+testutils.Mock1outof3pkgs(t))

		assert.Equal(t, "added 1 packages\n", actual)
	})

	t.Run("Db does not populate if items already exist", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "populate", "--dir="+t.TempDir())

		assert.Equal(t, "no new packages were found to populate the database, added 0 packages\n", actual)
	})
}
