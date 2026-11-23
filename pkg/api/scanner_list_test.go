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

// serveScanners points the client at a fake Harbor that holds total scanner
// registrations and pages /scanners like Harbor: page 1, page_size 10 by default.
func serveScanners(t *testing.T, total int) {
	t.Helper()

	scanners := make([]*models.ScannerRegistration, 0, total)
	for i := 1; i <= total; i++ {
		scanners = append(scanners, &models.ScannerRegistration{
			UUID: fmt.Sprintf("uuid-%d", i),
			Name: fmt.Sprintf("scanner-%d", i),
		})
	}

	queryInt := func(r *http.Request, key string, fallback int) int {
		if n, err := strconv.Atoi(r.URL.Query().Get(key)); err == nil {
			return n
		}
		return fallback
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2.0/scanners" {
			http.NotFound(w, r)
			return
		}
		page, pageSize := queryInt(r, "page", 1), queryInt(r, "page_size", 10)
		from := min((page-1)*pageSize, total)
		to := min(from+pageSize, total)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", strconv.Itoa(total))
		_ = json.NewEncoder(w).Encode(scanners[from:to])
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

func TestListScannersReturnsEveryPage(t *testing.T) {
	serveScanners(t, 215)

	resp, err := ListScanners()
	require.NoError(t, err)
	require.Len(t, resp.Payload, 215)
	assert.Equal(t, "scanner-1", resp.Payload[0].Name)
	assert.Equal(t, "scanner-215", resp.Payload[214].Name)
	assert.Equal(t, int64(215), resp.XTotalCount)
}

func TestListScannersFullLastPage(t *testing.T) {
	serveScanners(t, 200)

	resp, err := ListScanners()
	require.NoError(t, err)
	require.Len(t, resp.Payload, 200)
}

func TestGetScannerByNameBeyondFirstPage(t *testing.T) {
	serveScanners(t, 12)

	reg, err := GetScannerByName("scanner-12")
	require.NoError(t, err)
	assert.Equal(t, "uuid-12", reg.UUID)

	_, err = GetScannerByName("scanner-13")
	require.Error(t, err)
}
