package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/pkg/database"
	"github.com/vague2k/rummage/testutils"
)

func TestQuery(t *testing.T) {
	t.Run("Default amount of sorted matches", func(t *testing.T) {
		db, ctx := testutils.InMemDb(t)

		for i := range 10 {
			if i == 0 {
				continue
			}
			fakePkg := fmt.Sprintf("github.com/user%d/mux", i)

			_, err := db.AddItem(ctx, database.AddItemParams{
				Entry: fakePkg,
			})
			assert.NoError(t, err)

			err = db.UpdateItem(ctx, database.UpdateItemParams{
				Entry:        fakePkg,
				Score:        float64(i),
				Lastaccessed: int64(i),
			})
			assert.NoError(t, err)
		}

		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "mux")

		s := strings.Builder{}
		for i := 9; i > 0; i-- {
			s.WriteString(fmt.Sprintf("%d : %.4f : github.com/user%d/mux\n", int64(i), float64(i), i))
		}
		expected := s.String() + "\n"

		assert.Equal(t, expected, actual)
	})

	t.Run("Errors if no match was found", func(t *testing.T) {
		db, _ := testutils.InMemDb(t)
		cmd := NewRootCmd(db)
		actual := testutils.Execute(cmd, "query", "mux")

		assert.Equal(t, "no match found with the given arguement mux\n", actual)
	})
}
