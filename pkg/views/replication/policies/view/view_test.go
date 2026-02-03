package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatCopyByChunk(t *testing.T) {
	t.Run("nil CopyByChunk", func(t *testing.T) {
		got := formatCopyByChunk(nil)
		assert.Equal(t, "N/A", got)
	})

	t.Run("non-nil CopyByChunk", func(t *testing.T) {
		v := true
		got := formatCopyByChunk(&v)
		assert.Equal(t, "true", got)
	})
}
