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

//go:build e2e

package e2e

import (
	"net/http"
	"testing"
	"time"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/cveallowlist/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const systemCVEAllowlistPath = "/api/v2.0/system/CVEAllowlist"

// cveAllowlist is the subset of Harbor's allowlist payload the tests need.
type cveAllowlist struct {
	ExpiresAt *int64 `json:"expires_at"`
	Items     []struct {
		CVEID string `json:"cve_id"`
	} `json:"items"`
}

// preserveSystemCVEAllowlist restores the system allowlist when the test ends,
// since it is a single shared resource on the Harbor instance.
func preserveSystemCVEAllowlist(t *testing.T) {
	t.Helper()

	var original map[string]any
	doJSON(t, http.MethodGet, systemCVEAllowlistPath, nil, &original)
	t.Cleanup(func() {
		doJSON(t, http.MethodPut, systemCVEAllowlistPath, original, nil)
	})
}

func Test_CVEAllowlist_Add_ExpireDate(t *testing.T) {
	const cveID = "CVE-2026-1142"

	tests := []struct {
		name       string
		expireDate string
	}{
		{name: "documented dash format", expireDate: "2028-02-29"},
		{name: "legacy slash format", expireDate: "2028/02/29"},
	}

	want := time.Date(2028, time.February, 29, 0, 0, 0, 0, time.UTC).Unix()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginAndReload(t)
			preserveSystemCVEAllowlist(t)

			err := api.UpdateSystemCve(update.UpdateView{
				CveId:      cveID,
				IsExpire:   true,
				ExpireDate: tt.expireDate,
			})
			require.NoError(t, err)

			var got cveAllowlist
			doJSON(t, http.MethodGet, systemCVEAllowlistPath, nil, &got)

			require.NotNil(t, got.ExpiresAt)
			assert.Equal(t, want, *got.ExpiresAt)
			require.Len(t, got.Items, 1)
			assert.Equal(t, cveID, got.Items[0].CVEID)
		})
	}
}
