package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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

func selectItem(r *database.RummageDB, s string) *database.RummageDBItem {
	found, _ := r.SelectItem(s)
	return found
}

func resetScore(t *testing.T, r *database.RummageDB, item *database.RummageDBItem) {
	reset := &database.RummageDBItem{
		Entry: item.Entry,
		Score: 1.0,
	}
	_, err := r.UpdateItem(item.Entry, reset)
	assert.NoError(t, err)
}

// Mock the $GOPATH/pkg/mod/github.com dir using the test's TempDir
// All dirs (amt 100) should be "valid"
func mock100outof100pkgs(t *testing.T) string {
	dir := filepath.Join(t.TempDir(), "go", "pkg", "mod", "github.com")
	for i := range 100 {
		parentDir := filepath.Join(dir, fmt.Sprintf("dir%d", i))
		err := os.MkdirAll(parentDir, os.ModePerm)
		assert.NoError(t, err)

		childDir := filepath.Join(parentDir, fmt.Sprintf("child@%d.0.0", i))
		err = os.MkdirAll(childDir, os.ModePerm)
		assert.NoError(t, err)
	}

	return dir
}

// Mock the $GOPATH/pkg/mod/github.com dir using the test's TempDir
// Only 50 out of 100 dirs should be "valid"
func mock50outof100pkgs(t *testing.T) string {
	dir := filepath.Join(t.TempDir(), "go", "pkg", "mod", "github.com")
	for i := range 100 {
		parentDir := filepath.Join(dir, fmt.Sprintf("dir%d", i))
		if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
			t.Error(err)
		}

		if i > 49 {
			childDir := filepath.Join(parentDir, fmt.Sprintf("child@%d.0.0", i))
			file, err := os.Create(childDir)
			assert.NoError(t, err)
			file.Close()
		} else {
			childDir := filepath.Join(parentDir, fmt.Sprintf("child@%d.0.0", i))
			err := os.MkdirAll(childDir, os.ModePerm)
			assert.NoError(t, err)
		}
	}

	return dir
}
