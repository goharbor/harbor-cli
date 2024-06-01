package quota

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageStringToBytes(t *testing.T) {
	// Valid inputs
	tests := []struct {
		input    string
		expected int64
	}{
		{"1MiB", 1024 * 1024},
		{"1GiB", 1024 * 1024 * 1024},
		{"1TiB", 1024 * 1024 * 1024 * 1024},
	}

	for _, test := range tests {
		result, err := storageStringToBytes(test.input)
		assert.NoError(t, err, "Unexpected error for input %s", test.input)
		assert.Equal(
			t,
			test.expected,
			result,
			"Expected %d but got %d for input %s",
			test.expected,
			result,
			test.input,
		)
	}

	// Invalid inputs
	invalidInputs := []string{
		"1KB",
		"1000",
		"10PB",
		"1GiBGiB",
		"1.03GiB",
		"1.08TiB",
	}

	for _, input := range invalidInputs {
		_, err := storageStringToBytes(input)
		assert.Error(t, err, "Expected error for input %s but got none", input)
	}

	// Exceeding maximum value
	_, err := storageStringToBytes("1025TiB")
	assert.Error(t, err, "Expected error for input exceeding 1024TiB but got none")
}
