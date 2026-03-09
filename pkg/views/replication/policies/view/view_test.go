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

	t.Run("true CopyByChunk", func(t *testing.T) {
		v := true
		got := formatCopyByChunk(&v)
		assert.Equal(t, "true", got)
	})

	t.Run("false CopyByChunk", func(t *testing.T) {
		v := false
		got := formatCopyByChunk(&v)
		assert.Equal(t, "false", got)
	})
}

func TestFormatSpeed(t *testing.T) {
	t.Run("nil Speed", func(t *testing.T) {
		got := formatSpeed(nil)
		assert.Equal(t, "N/A", got)
	})

	t.Run("non-nil Speed", func(t *testing.T) {
		v := int32(1024)
		got := formatSpeed(&v)
		assert.Equal(t, "1024 B/s", got)
	})
}
