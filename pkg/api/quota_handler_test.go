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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/quota"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllQuotasFetchesAllPages(t *testing.T) {
	var pages []int64
	var pageSizes []int64

	quotas, err := GetAllQuotas(
		func(opts ListQuotaFlags) (*quota.ListQuotasOK, error) {
			pages = append(pages, opts.Page)
			pageSizes = append(pageSizes, opts.PageSize)

			if opts.Page == 1 {
				return &quota.ListQuotasOK{Payload: makeQuotas(100)}, nil
			}
			return &quota.ListQuotasOK{Payload: makeQuotas(2)}, nil
		},
		ListQuotaFlags{Page: 8, PageSize: 0},
	)

	require.NoError(t, err)
	assert.Len(t, quotas, 102)
	assert.Equal(t, []int64{1, 2}, pages)
	assert.Equal(t, []int64{100, 100}, pageSizes)
}

func makeQuotas(count int) []*models.Quota {
	quotas := make([]*models.Quota, count)
	for i := range quotas {
		quotas[i] = &models.Quota{ID: int64(i + 1)}
	}
	return quotas
}
