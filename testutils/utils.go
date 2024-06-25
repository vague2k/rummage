package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/vague2k/rummage/pkg/database"
)

func DbInstance(t *testing.T) *database.RummageDB {
	tmp := t.TempDir()
	r, err := database.Init(tmp)
	if err != nil {
		t.Errorf("Could not open db: \n%s", err)
	}
	return r
}

func AssertEquals(t *testing.T, expected any, got any) {
	if expected != got {
		t.Errorf("Expected %v, but got %v.", expected, got)
	}
	t.Log(expected, got)
}

func AssertNotEquals(t *testing.T, expected any, got any) {
	if expected == got {
		t.Errorf("Expected %v, but got %v.", expected, got)
	}
	t.Log(expected, got)
}

func AssertNotNil(t *testing.T, got any) {
	if got == nil {
		t.Errorf("Expected %v, but got %v.", nil, got)
	}
	t.Log(nil, got)
}

func AssertTrue(t *testing.T, got bool) {
	if !got {
		t.Errorf("Expected %v, but got %v.", true, got)
	}
	t.Log(true, got)
}

func AssertFalse(t *testing.T, got bool) {
	if got {
		t.Errorf("Expected %v, but got %v.", false, got)
	}
	t.Log(false, got)
}

func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func MockPkgModGHDir(dir string) error {
	for i := range 10 {
		parentDir := filepath.Join(dir, fmt.Sprintf("dir%d", i))
		if err := os.Mkdir(parentDir, os.ModePerm); err != nil {
			return err
		}

		childDir := filepath.Join(parentDir, fmt.Sprintf("child@%d.0.0", i))
		if err := os.Mkdir(childDir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
