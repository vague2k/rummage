package services_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/vague2k/rummage/pkg/commands"
	"github.com/vague2k/rummage/testutils"
)

func TestWalkAndParsePackages(t *testing.T) {
	t.Run("Returns 100 valid packages, out of 100", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "go", "pkg", "mod", "github.com")
		// quickly mock the $GOPATH/pkg/mod/github.com dir using the test's TempDir
		// All dirs (amt 100) should be "valid"
		for i := range 100 {
			parentDir := filepath.Join(dir, fmt.Sprintf("dir%d", i))
			if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
				t.Error(err)
			}

			childDir := filepath.Join(parentDir, fmt.Sprintf("child@%d.0.0", i))
			if err := os.MkdirAll(childDir, os.ModePerm); err != nil {
				t.Error(err)
			}
		}

		pkgs := commands.WalkAndParsePackages(dir)

		testutils.AssertEquals(t, 100, len(pkgs))
		for _, pkg := range pkgs {
			t.Log(pkg)
		}
	})

	t.Run("Returns 50 valid packages, out of 100", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "go", "pkg", "mod", "github.com")
		// quickly mock the $GOPATH/pkg/mod/github.com dir using the test's TempDir
		// Only 50 out of 100 dirs should be "valid"
		for i := range 100 {
			parentDir := filepath.Join(dir, fmt.Sprintf("dir%d", i))
			if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
				t.Error(err)
			}

			if i > 49 {
				childDir := filepath.Join(parentDir, fmt.Sprintf("child@%d.0.0", i))
				file, err := os.Create(childDir)
				testutils.CheckErr(t, err)
				file.Close()
			} else {
				childDir := filepath.Join(parentDir, fmt.Sprintf("child@%d.0.0", i))
				if err := os.MkdirAll(childDir, os.ModePerm); err != nil {
					t.Error(err)
				}
			}
		}
		pkgs := commands.WalkAndParsePackages(dir)

		testutils.AssertEquals(t, 50, len(pkgs))
		for _, pkg := range pkgs {
			t.Log(pkg)
		}
	})
}
