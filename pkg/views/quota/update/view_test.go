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
package update

import (
	"strconv"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageDefaults(t *testing.T) {
	tests := []struct {
		name      string
		bytes     int64
		wantValue string
		wantUnit  string
	}{
		{"whole gibibytes", 50 * gibibyte, "50", "GiB"},
		{"whole tebibytes", 2 * tebibyte, "2", "TiB"},
		{"whole mebibytes", 512 * mebibyte, "512", "MiB"},
		{"falls back to mebibytes when not a whole gibibyte", 1536 * mebibyte, "1536", "MiB"},
		{"zero limit", 0, "0", "MiB"},
		{"unlimited limit has no prefilled value", unlimitedStorage, "", defaultStorageUnit},
		{"limit smaller than a mebibyte has no prefilled value", 1024, "", defaultStorageUnit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, unit := storageDefaults(tt.bytes)
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantUnit, unit)
		})
	}
}

// The quota limit input rejects anything ParseInt cannot read, and the storage
// unit select only offers the units below, so both defaults must be accepted by
// the form that receives them.
func TestStorageDefaultsAreAcceptedByTheForm(t *testing.T) {
	limits := []int64{0, mebibyte, 512 * mebibyte, gibibyte, 50 * gibibyte, 2 * tebibyte}

	for _, limit := range limits {
		t.Run(strconv.FormatInt(limit, 10), func(t *testing.T) {
			value, unit := storageDefaults(limit)

			require.NotEmpty(t, value)
			_, err := strconv.ParseInt(value, 10, 64)
			assert.NoError(t, err)
			assert.Contains(t, storageUnits, unit)
		})
	}
}

// The view returns value and unit concatenated, and the caller feeds that back
// into StorageStringToBytes, so an untouched form must round-trip to the limit
// it was prefilled with.
func TestStorageDefaultsRoundTrip(t *testing.T) {
	limits := []int64{mebibyte, 512 * mebibyte, gibibyte, 50 * gibibyte, 2 * tebibyte}

	for _, limit := range limits {
		t.Run(strconv.FormatInt(limit, 10), func(t *testing.T) {
			value, unit := storageDefaults(limit)

			got, err := utils.StorageStringToBytes(value + unit)
			require.NoError(t, err)
			assert.Equal(t, limit, got)
		})
	}
}

func TestCurrentStorage(t *testing.T) {
	t.Run("unlimited quota", func(t *testing.T) {
		assert.Equal(t, "Unlimited", currentStorage(unlimitedStorage))
	})

	t.Run("limited quota", func(t *testing.T) {
		assert.Equal(t, "50.0 GiB", currentStorage(50*gibibyte))
	})
}
