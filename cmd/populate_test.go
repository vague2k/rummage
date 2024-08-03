package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vague2k/rummage/mocks"
	"github.com/vague2k/rummage/pkg/database"
)

func TestPopulate(t *testing.T) {
	t.Run("Can populate db with 3 out of 3 packages", func(t *testing.T) {
		mockDb := &mocks.RummageDbInterface{}
		returnVal := []*database.RummageItem{{}, {}, {}} // actual item values are not needed
		mockDb.On("AddMultiItems", "github.com/dir0/child", "github.com/dir1/child", "github.com/dir2/child").Return(returnVal, 3, nil).Once()

		cmd := NewRootCmd(mockDb)
		actual := execute(cmd, "populate", "--dir="+mock3outof3pkgs(t))

		assert.Equal(t, "added 3 packages\n", actual)
		mockDb.AssertExpectations(t)
	})

	t.Run("Can populate db with 1 out of 3 packages", func(t *testing.T) {
		mockDb := &mocks.RummageDbInterface{}
		returnVal := []*database.RummageItem{{}} // actual item values are not needed
		mockDb.On("AddMultiItems", "github.com/dir0/child").Return(returnVal, 1, nil).Once()

		cmd := NewRootCmd(mockDb)
		actual := execute(cmd, "populate", "--dir="+mock1outof3pkgs(t))

		assert.Equal(t, "added 1 packages\n", actual)
		mockDb.AssertExpectations(t)
	})

	t.Run("Db does not populate if items already exist", func(t *testing.T) {
		mockDb := &mocks.RummageDbInterface{}
		mockDb.On("AddMultiItems", "github.com/dir0/child").Return(nil, 0, nil).Once()
		cmd := NewRootCmd(mockDb)
		actual := execute(cmd, "populate", "--dir="+mock1outof3pkgs(t))

		assert.Equal(t, "no new packages were found to populate the database, added 0 packages\n", actual)
		mockDb.AssertExpectations(t)
	})
}
