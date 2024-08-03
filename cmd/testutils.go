package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func execute(cmd *cobra.Command, args ...string) string {
	buf := new(bytes.Buffer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	cmd.Execute()

	return buf.String()
}

// Mock the $GOPATH/pkg/mod/github.com dir using the test's TempDir
// Only 3 out of 3 dirs should be "valid"
func mock3outof3pkgs(t *testing.T) string {
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
// Only 1 out of 3 dirs should be "valid"
func mock1outof3pkgs(t *testing.T) string {
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
