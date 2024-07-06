package services_test

import (
	"testing"

	s "github.com/vague2k/rummage/cmd/services"
	"github.com/vague2k/rummage/testutils"
)

func TestAttemptGoGet(t *testing.T) {
	t.Run("Successfully gets go package", func(t *testing.T) {
		test := "github.com/gorilla/mux"
		err := s.AttemptGoGet(test)
		if err != nil {
			t.Errorf("Expected no error, but got an error: \n%s", err)
		}
	})

	t.Run("Errors when not a go package", func(t *testing.T) {
		test := "notagopackage"
		err := s.AttemptGoGet(test)
		if err == nil {
			t.Errorf("Expected error, but got nil: \n%s", err)
		}
	})
}

func TestUpdateRecency(t *testing.T) {
	r := testutils.DbInstance(t)
	defer r.DB.Close()

	item, err := r.AddItem("item")
	if err != nil {
		t.Errorf("Issue adding item to db: \n%s", err)
	}
	t.Run("Properly increments score", func(t *testing.T) {
		got := s.UpdateRecency(r, item)
		expected := 4.0

		testutils.AssertEquals(t, expected, got.Score)
	})
}
