package testutils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/database"
)

// Spin up an in memory database (since we're using sqlite3) for quick testing
//
// This function already includes a cleanup function where when the test completes, the database is closed
func InMemDb(t *testing.T) *database.RummageDb {
	db, err := database.Init(":memory:")
	assert.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
		db = nil
	})
	return db
}

// Execute a cobra command with custom args for testing using *bytes.Buffer under the hood.
//
// Since the root cmd is used for testing, any subcommands the root command has will have to be included in the args
func Execute(cmd *cobra.Command, args ...string) string {
	buf := new(bytes.Buffer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	return buf.String()
}

// Mock the $GOPATH/pkg/mod/github.com dir using the test's TempDir
//
// Only 3 out of 3 dirs should be "valid"
func Mock3outof3pkgs(t *testing.T) string {
	dir := filepath.Join(t.TempDir(), "go", "pkg", "mod", "github.com")
	for i := range 3 {
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
//
// Only 1 out of 3 dirs should be "valid"
func Mock1outof3pkgs(t *testing.T) string {
	dir := filepath.Join(t.TempDir(), "go", "pkg", "mod", "github.com")
	for i := range 3 {
		parentDir := filepath.Join(dir, fmt.Sprintf("dir%d", i))
		if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
			t.Error(err)
		}

		if i == 0 {
			childDir := filepath.Join(parentDir, fmt.Sprintf("child@%d.0.0", i))
			err := os.MkdirAll(childDir, os.ModePerm)
			assert.NoError(t, err)
		} else {
			childDir := filepath.Join(parentDir, fmt.Sprintf("file@%d.0.0", i))
			file, err := os.Create(childDir)
			assert.NoError(t, err)
			file.Close()
		}
	}

	return dir
}

// Use this for the "get" command tests.
//
// When running these tests locally it can affect our go.mod file,
// which would need to get cleaned up after every test
//
// This is not the most clever way of handling this issue, but it's a whole
// FUCK of alot easier then... idk say spinning up a docker container with
// a dedicated go project directory just for these tests. And it suprisingly
// works a lot better than I thought it would
func GoModTidy(t *testing.T) {
	cmd := exec.Command("go", "mod", "tidy")
	err := cmd.Run()
	assert.NoError(t, err)
}
