package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulate(t *testing.T) {
	t.Run("Returns 100 valid packages, out of 100", func(t *testing.T) {
		dir := mock100outof100pkgs(t)
		pkgs, err := extractPackages(dir)
		assert.NoError(t, err)
		assert.Len(t, pkgs, 100)
	})

	t.Run("Returns 50 valid packages, out of 100", func(t *testing.T) {
		dir := mock50outof100pkgs(t)
		pkgs, err := extractPackages(dir)
		assert.NoError(t, err)
		assert.Len(t, pkgs, 50)
	})
}
