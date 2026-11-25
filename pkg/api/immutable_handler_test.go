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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// usePagedImmutableServer serves total immutable rules the way Harbor does:
// page_size defaults to 10 and page to 1.
func usePagedImmutableServer(t *testing.T, total int) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := 1, 10
		if v := r.URL.Query().Get("page"); v != "" {
			page, _ = strconv.Atoi(v)
		}
		if v := r.URL.Query().Get("page_size"); v != "" {
			pageSize, _ = strconv.Atoi(v)
		}

		rules := []*models.ImmutableRule{}
		for i := (page - 1) * pageSize; i < total && i < page*pageSize; i++ {
			rules = append(rules, &models.ImmutableRule{ID: int64(i + 1)})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rules)
	}))

	utils.ClientOnce = sync.Once{}
	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = utils.GetClientByConfig(&harbor.ClientSetConfig{URL: srv.URL})
	utils.ClientErr = nil

	t.Cleanup(func() {
		srv.Close()
		utils.ClientOnce = sync.Once{}
		utils.ClientInstance = nil
		utils.ClientErr = nil
	})
}

func TestListImmutableReadsAllPages(t *testing.T) {
	for _, total := range []int{0, 3, 25, 100, 250} {
		t.Run(strconv.Itoa(total), func(t *testing.T) {
			usePagedImmutableServer(t, total)

			resp, err := ListImmutable("library")
			require.NoError(t, err)
			require.Len(t, resp.Payload, total)
			for i, rule := range resp.Payload {
				assert.Equal(t, int64(i+1), rule.ID)
			}
		})
	}
}
