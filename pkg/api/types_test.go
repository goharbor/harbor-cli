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
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestListFlags_JSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ListFlags
		expected ListFlags
	}{
		{
			name: "all fields populated",
			input: ListFlags{
				ProjectID: 123,
				Scope:     "project",
				Name:      "test-name",
				Page:      1,
				PageSize:  10,
				Q:         "query",
				Sort:      "name",
				Public:    true,
			},
			expected: ListFlags{
				ProjectID: 123,
				Scope:     "project",
				Name:      "test-name",
				Page:      1,
				PageSize:  10,
				Q:         "query",
				Sort:      "name",
				Public:    true,
			},
		},
		{
			name: "zero values",
			input: ListFlags{
				ProjectID: 0,
				Scope:     "",
				Name:      "",
				Page:      0,
				PageSize:  0,
				Q:         "",
				Sort:      "",
				Public:    false,
			},
			expected: ListFlags{
				ProjectID: 0,
				Scope:     "",
				Name:      "",
				Page:      0,
				PageSize:  0,
				Q:         "",
				Sort:      "",
				Public:    false,
			},
		},
		{
			name: "partial fields",
			input: ListFlags{
				ProjectID: 456,
				Page:      2,
				PageSize:  20,
				Public:    false,
			},
			expected: ListFlags{
				ProjectID: 456,
				Scope:     "",
				Name:      "",
				Page:      2,
				PageSize:  20,
				Q:         "",
				Sort:      "",
				Public:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			var result ListFlags
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateRegView_JSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    CreateRegView
		expected CreateRegView
	}{
		{
			name: "all fields populated",
			input: CreateRegView{
				Name:        "test-registry",
				Type:        "harbor",
				Description: "Test registry description",
				URL:         "https://registry.example.com",
				Credential: RegistryCredential{
					AccessKey:    "access-key",
					Type:         "basic",
					AccessSecret: "secret",
				},
				Insecure: true,
			},
			expected: CreateRegView{
				Name:        "test-registry",
				Type:        "harbor",
				Description: "Test registry description",
				URL:         "https://registry.example.com",
				Credential: RegistryCredential{
					AccessKey:    "access-key",
					Type:         "basic",
					AccessSecret: "secret",
				},
				Insecure: true,
			},
		},
		{
			name: "zero values",
			input: CreateRegView{
				Name:        "",
				Type:        "",
				Description: "",
				URL:         "",
				Credential:  RegistryCredential{},
				Insecure:    false,
			},
			expected: CreateRegView{
				Name:        "",
				Type:        "",
				Description: "",
				URL:         "",
				Credential:  RegistryCredential{},
				Insecure:    false,
			},
		},
		{
			name: "empty credential",
			input: CreateRegView{
				Name:        "registry-name",
				Type:        "docker-hub",
				Description: "Description",
				URL:         "https://hub.docker.com",
				Credential:  RegistryCredential{},
				Insecure:    false,
			},
			expected: CreateRegView{
				Name:        "registry-name",
				Type:        "docker-hub",
				Description: "Description",
				URL:         "https://hub.docker.com",
				Credential:  RegistryCredential{},
				Insecure:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			var result CreateRegView
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContextListView_JSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ContextListView
		expected ContextListView
	}{
		{
			name: "all fields populated",
			input: ContextListView{
				Name:     "test-context",
				Username: "test-user",
				Server:   "https://harbor.example.com",
			},
			expected: ContextListView{
				Name:     "test-context",
				Username: "test-user",
				Server:   "https://harbor.example.com",
			},
		},
		{
			name: "zero values",
			input: ContextListView{
				Name:     "",
				Username: "",
				Server:   "",
			},
			expected: ContextListView{
				Name:     "",
				Username: "",
				Server:   "",
			},
		},
		{
			name: "partial fields",
			input: ContextListView{
				Name:   "context-name",
				Server: "https://server.com",
			},
			expected: ContextListView{
				Name:     "context-name",
				Username: "",
				Server:   "https://server.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			var result ContextListView
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRegistryCredential_JSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    RegistryCredential
		expected RegistryCredential
		jsonStr  string
	}{
		{
			name: "all fields populated",
			input: RegistryCredential{
				AccessKey:    "access-key-123",
				Type:         "basic",
				AccessSecret: "secret-456",
			},
			expected: RegistryCredential{
				AccessKey:    "access-key-123",
				Type:         "basic",
				AccessSecret: "secret-456",
			},
		},
		{
			name: "zero values",
			input: RegistryCredential{
				AccessKey:    "",
				Type:         "",
				AccessSecret: "",
			},
			expected: RegistryCredential{
				AccessKey:    "",
				Type:         "",
				AccessSecret: "",
			},
		},
		{
			name: "empty fields with omitempty",
			input: RegistryCredential{},
			expected: RegistryCredential{},
			jsonStr:  "{}",
		},
		{
			name: "partial fields",
			input: RegistryCredential{
				AccessKey: "key-only",
			},
			expected: RegistryCredential{
				AccessKey: "key-only",
			},
		},
		{
			name: "JSON with snake_case tags",
			input: RegistryCredential{
				AccessKey:    "test-key",
				Type:         "oauth",
				AccessSecret: "test-secret",
			},
			expected: RegistryCredential{
				AccessKey:    "test-key",
				Type:         "oauth",
				AccessSecret: "test-secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Marshal
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			// If jsonStr is provided, verify it matches
			if tt.jsonStr != "" {
				assert.JSONEq(t, tt.jsonStr, string(data))
			}

			// Test Unmarshal
			var result RegistryCredential
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)

			// Test unmarshaling from JSON with snake_case
			if tt.name == "JSON with snake_case tags" {
				jsonData := `{"access_key":"test-key","type":"oauth","access_secret":"test-secret"}`
				var snakeResult RegistryCredential
				err = json.Unmarshal([]byte(jsonData), &snakeResult)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, snakeResult)
			}
		})
	}
}

func TestListQuotaFlags_JSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ListQuotaFlags
		expected ListQuotaFlags
	}{
		{
			name: "all fields populated",
			input: ListQuotaFlags{
				PageSize:    50,
				Page:        3,
				Sort:        "hard",
				Reference:   "project",
				ReferenceID: "ref-123",
			},
			expected: ListQuotaFlags{
				PageSize:    50,
				Page:        3,
				Sort:        "hard",
				Reference:   "project",
				ReferenceID: "ref-123",
			},
		},
		{
			name: "zero values",
			input: ListQuotaFlags{
				PageSize:    0,
				Page:        0,
				Sort:        "",
				Reference:   "",
				ReferenceID: "",
			},
			expected: ListQuotaFlags{
				PageSize:    0,
				Page:        0,
				Sort:        "",
				Reference:   "",
				ReferenceID: "",
			},
		},
		{
			name: "partial fields",
			input: ListQuotaFlags{
				Page:      1,
				PageSize:  25,
				Reference: "project",
			},
			expected: ListQuotaFlags{
				PageSize:    25,
				Page:        1,
				Sort:        "",
				Reference:   "project",
				ReferenceID: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			var result ListQuotaFlags
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestListMemberOptions_JSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ListMemberOptions
		expected ListMemberOptions
	}{
		{
			name: "all fields populated",
			input: ListMemberOptions{
				XIsResourceName: true,
				ProjectNameOrID:  "project-123",
				Page:             2,
				PageSize:         30,
				EntityName:       "entity-name",
				WithDetail:       true,
			},
			expected: ListMemberOptions{
				XIsResourceName: true,
				ProjectNameOrID:  "project-123",
				Page:             2,
				PageSize:         30,
				EntityName:       "entity-name",
				WithDetail:       true,
			},
		},
		{
			name: "zero values",
			input: ListMemberOptions{
				XIsResourceName: false,
				ProjectNameOrID:  "",
				Page:             0,
				PageSize:         0,
				EntityName:       "",
				WithDetail:       false,
			},
			expected: ListMemberOptions{
				XIsResourceName: false,
				ProjectNameOrID:  "",
				Page:             0,
				PageSize:         0,
				EntityName:       "",
				WithDetail:       false,
			},
		},
		{
			name: "partial fields",
			input: ListMemberOptions{
				ProjectNameOrID: "my-project",
				Page:            1,
				PageSize:        10,
			},
			expected: ListMemberOptions{
				XIsResourceName: false,
				ProjectNameOrID:  "my-project",
				Page:             1,
				PageSize:         10,
				EntityName:       "",
				WithDetail:       false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			var result ListMemberOptions
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUpdateMemberOptions_JSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    UpdateMemberOptions
		expected UpdateMemberOptions
	}{
		{
			name: "all fields populated with RoleID",
			input: UpdateMemberOptions{
				XIsResourceName: true,
				ID:              42,
				ProjectNameOrID:  "project-456",
				RoleID: &models.RoleRequest{
					RoleID: 1,
				},
			},
			expected: UpdateMemberOptions{
				XIsResourceName: true,
				ID:              42,
				ProjectNameOrID:  "project-456",
				RoleID: &models.RoleRequest{
					RoleID: 1,
				},
			},
		},
		{
			name: "zero values with nil RoleID",
			input: UpdateMemberOptions{
				XIsResourceName: false,
				ID:              0,
				ProjectNameOrID:  "",
				RoleID:          nil,
			},
			expected: UpdateMemberOptions{
				XIsResourceName: false,
				ID:              0,
				ProjectNameOrID:  "",
				RoleID:          nil,
			},
		},
		{
			name: "partial fields with RoleID",
			input: UpdateMemberOptions{
				ID:             100,
				ProjectNameOrID: "test-project",
				RoleID: &models.RoleRequest{
					RoleID: 2,
				},
			},
			expected: UpdateMemberOptions{
				XIsResourceName: false,
				ID:              100,
				ProjectNameOrID:  "test-project",
				RoleID: &models.RoleRequest{
					RoleID: 2,
				},
			},
		},
		{
			name: "RoleID with zero value",
			input: UpdateMemberOptions{
				ID:             50,
				ProjectNameOrID: "project",
				RoleID: &models.RoleRequest{
					RoleID: 0,
				},
			},
			expected: UpdateMemberOptions{
				XIsResourceName: false,
				ID:              50,
				ProjectNameOrID:  "project",
				RoleID: &models.RoleRequest{
					RoleID: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			var result UpdateMemberOptions
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.XIsResourceName, result.XIsResourceName)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.ProjectNameOrID, result.ProjectNameOrID)
			if tt.expected.RoleID == nil {
				assert.Nil(t, result.RoleID)
			} else {
				assert.NotNil(t, result.RoleID)
				assert.Equal(t, tt.expected.RoleID.RoleID, result.RoleID.RoleID)
			}
		})
	}
}

func TestGetMemberOptions_JSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    GetMemberOptions
		expected GetMemberOptions
	}{
		{
			name: "all fields populated",
			input: GetMemberOptions{
				XIsResourceName: true,
				ID:              789,
				ProjectNameOrID:  "project-789",
			},
			expected: GetMemberOptions{
				XIsResourceName: true,
				ID:              789,
				ProjectNameOrID:  "project-789",
			},
		},
		{
			name: "zero values",
			input: GetMemberOptions{
				XIsResourceName: false,
				ID:              0,
				ProjectNameOrID:  "",
			},
			expected: GetMemberOptions{
				XIsResourceName: false,
				ID:              0,
				ProjectNameOrID:  "",
			},
		},
		{
			name: "partial fields",
			input: GetMemberOptions{
				ID:             999,
				ProjectNameOrID: "test-project-id",
			},
			expected: GetMemberOptions{
				XIsResourceName: false,
				ID:              999,
				ProjectNameOrID:  "test-project-id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			var result GetMemberOptions
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRegistryCredential_JSONTags(t *testing.T) {
	t.Run("marshal uses snake_case tags", func(t *testing.T) {
		cred := RegistryCredential{
			AccessKey:    "key",
			Type:         "type",
			AccessSecret: "secret",
		}

		data, err := json.Marshal(cred)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		assert.NoError(t, err)

		// Verify snake_case keys are present
		assert.Contains(t, result, "access_key")
		assert.Contains(t, result, "type")
		assert.Contains(t, result, "access_secret")
		assert.Equal(t, "key", result["access_key"])
		assert.Equal(t, "type", result["type"])
		assert.Equal(t, "secret", result["access_secret"])
	})

	t.Run("unmarshal from snake_case JSON", func(t *testing.T) {
		jsonData := `{"access_key":"test-key","type":"basic","access_secret":"test-secret"}`
		var cred RegistryCredential

		err := json.Unmarshal([]byte(jsonData), &cred)
		assert.NoError(t, err)
		assert.Equal(t, "test-key", cred.AccessKey)
		assert.Equal(t, "basic", cred.Type)
		assert.Equal(t, "test-secret", cred.AccessSecret)
	})

	t.Run("omitempty works for empty fields", func(t *testing.T) {
		cred := RegistryCredential{}

		data, err := json.Marshal(cred)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		assert.NoError(t, err)

		// With omitempty, empty fields should not appear in JSON
		assert.Empty(t, result)
	})
}
