// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

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
		result, err := StorageStringToBytes(test.input)
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
		_, err := StorageStringToBytes(input)
		assert.Error(t, err, "Expected error for input %s but got none", input)
	}

	// Exceeding maximum value
	_, err := StorageStringToBytes("1025TiB")
	assert.Error(t, err, "Expected error for input exceeding 1024TiB but got none")
}
