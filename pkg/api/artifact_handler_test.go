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

func TestListTagsFetchesAllPages(t *testing.T) {
	const total = 230
	var tags []*models.Tag
	for i := range total {
		tags = append(tags, &models.Tag{Name: fmt.Sprintf("v%d", i)})
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := 1, 10 // Harbor's defaults
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
		_ = json.NewEncoder(w).Encode(tags[start:end])
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

	resp, err := ListTags("library", "app", "latest")
	require.NoError(t, err)
	require.Len(t, resp.Payload, total)
	assert.Equal(t, "v0", resp.Payload[0].Name)
	assert.Equal(t, "v229", resp.Payload[total-1].Name)
}
