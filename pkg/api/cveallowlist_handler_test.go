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
package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseExpireDate(t *testing.T) {
	leapDay := time.Date(2028, time.February, 29, 0, 0, 0, 0, time.UTC).Unix()

	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{name: "dash format", input: "2028-02-29", want: leapDay},
		{name: "slash format", input: "2028/02/29", want: leapDay},
		{name: "invalid day", input: "2027-02-29", wantErr: true},
		{name: "mixed separators", input: "2028-02/29", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExpireDate(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
