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
	"fmt"
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

// newWebhookServer serves total webhook policies for project "library" and
// pages them the way Harbor does: page 1 and page_size 10 when not given.
func newWebhookServer(t *testing.T, total int) {
	t.Helper()

	var policies []*models.WebhookPolicy
	for i := range total {
		policies = append(policies, &models.WebhookPolicy{ID: int64(i + 1), Name: fmt.Sprintf("hook-%d", i+1)})
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2.0/projects/library/webhook/policies" {
			http.NotFound(w, r)
			return
		}
		page, pageSize := 1, 10
		if v := r.URL.Query().Get("page"); v != "" {
			page, _ = strconv.Atoi(v)
		}
		if v := r.URL.Query().Get("page_size"); v != "" {
			pageSize, _ = strconv.Atoi(v)
		}
		start := min((page-1)*pageSize, total)
		end := min(start+pageSize, total)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", strconv.Itoa(total))
		_ = json.NewEncoder(w).Encode(policies[start:end])
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

func TestListWebhooksFetchesAllPages(t *testing.T) {
	const total = 230
	newWebhookServer(t, total)

	resp, err := ListWebhooks("library")
	require.NoError(t, err)
	require.Len(t, resp.Payload, total)
	assert.Equal(t, "hook-1", resp.Payload[0].Name)
	assert.Equal(t, "hook-230", resp.Payload[total-1].Name)
	assert.Equal(t, int64(total), resp.XTotalCount)
}

func TestListWebhooksExactPageSize(t *testing.T) {
	newWebhookServer(t, 100)

	resp, err := ListWebhooks("library")
	require.NoError(t, err)
	require.Len(t, resp.Payload, 100)
}

func TestGetWebhookIDFindsPolicyPastFirstPage(t *testing.T) {
	newWebhookServer(t, 25)

	id, err := GetWebhookID("library", "hook-25")
	require.NoError(t, err)
	assert.Equal(t, int64(25), id)

	_, err = GetWebhookID("library", "missing")
	require.Error(t, err)
}
