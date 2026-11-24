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
package selection

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelected(t *testing.T) {
	tests := map[string]struct {
		model   tea.Model
		want    string
		wantErr error
	}{
		"chosen":       {model: Model{Choice: "library"}, want: "library"},
		"aborted":      {model: Model{Aborted: true}, wantErr: ErrUserAborted},
		"nothing":      {model: Model{}, wantErr: ErrNoSelection},
		"foreign type": {model: nil, wantErr: ErrNoSelection},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := Selected(tt.model)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
