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
package list

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookRows(t *testing.T) {
	t.Run("policy with a target", func(t *testing.T) {
		rows := webhookRows([]*models.WebhookPolicy{
			{
				ID:      1,
				Name:    "hook",
				Enabled: true,
				Targets: []*models.WebhookTargetObject{
					{
						Address:       "https://example.com/hook",
						Type:          "http",
						PayloadFormat: models.PayloadFormatType("CloudEvents"),
					},
				},
			},
		})

		require.Len(t, rows, 1)
		assert.Equal(t, "1", rows[0][0])
		assert.Equal(t, "hook", rows[0][1])
		assert.Equal(t, "True", rows[0][2])
		assert.Equal(t, "https://example.com/hook", rows[0][3])
		assert.Equal(t, "http", rows[0][4])
		assert.Equal(t, "CloudEvents", rows[0][5])
	})

	t.Run("policy without targets falls back to placeholders", func(t *testing.T) {
		rows := webhookRows([]*models.WebhookPolicy{
			{ID: 2, Name: "no-targets"},
		})

		require.Len(t, rows, 1)
		assert.Equal(t, "2", rows[0][0])
		assert.Equal(t, "no-targets", rows[0][1])
		assert.Equal(t, "False", rows[0][2])
		assert.Equal(t, placeholder, rows[0][3])
		assert.Equal(t, placeholder, rows[0][4])
		assert.Equal(t, placeholder, rows[0][5])
	})

	t.Run("policy with a nil target falls back to placeholders", func(t *testing.T) {
		rows := webhookRows([]*models.WebhookPolicy{
			{ID: 3, Name: "nil-target", Targets: []*models.WebhookTargetObject{nil}},
		})

		require.Len(t, rows, 1)
		assert.Equal(t, placeholder, rows[0][3])
		assert.Equal(t, placeholder, rows[0][4])
		assert.Equal(t, placeholder, rows[0][5])
	})

	t.Run("target with empty fields falls back to placeholders", func(t *testing.T) {
		rows := webhookRows([]*models.WebhookPolicy{
			{ID: 4, Name: "empty-target", Targets: []*models.WebhookTargetObject{{}}},
		})

		require.Len(t, rows, 1)
		assert.Equal(t, placeholder, rows[0][3])
		assert.Equal(t, placeholder, rows[0][4])
		assert.Equal(t, placeholder, rows[0][5])
	})

	t.Run("nil policies are skipped", func(t *testing.T) {
		rows := webhookRows([]*models.WebhookPolicy{nil, {ID: 5, Name: "kept"}})

		require.Len(t, rows, 1)
		assert.Equal(t, "kept", rows[0][1])
	})

	t.Run("row has one cell per column", func(t *testing.T) {
		rows := webhookRows([]*models.WebhookPolicy{{ID: 6}})

		require.Len(t, rows, 1)
		assert.Len(t, rows[0], len(columns))
	})

	t.Run("no policies yields no rows", func(t *testing.T) {
		assert.Empty(t, webhookRows(nil))
	})
}
