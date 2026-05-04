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
package project

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestSummaryCommand_Success(t *testing.T) {
	originalFunc := getProjectSummaryFunc
	defer func() { getProjectSummaryFunc = originalFunc }()

	getProjectSummaryFunc = func(projectNameOrID string, useProjectID bool) (*project.GetProjectSummaryOK, error) {
		return &project.GetProjectSummaryOK{
			Payload: &models.ProjectSummary{
				RepoCount:         1,
				ProjectAdminCount: 1,
				Quota: &models.ProjectSummaryQuota{
					Hard: models.ResourceList{"storage": 1024 * 1024 * 1024},
					Used: models.ResourceList{"storage": 512},
				},
				Registry: &models.Registry{
					Name:   "test-registry",
					URL:    "https://registry.test",
					Type:   "docker-hub",
					Status: "healthy",
				},
			},
		}, nil
	}

	cmd := SummaryCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"test-project"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestSummaryCommand_Errors(t *testing.T) {
	tests := []struct {
		name          string
		mockErr       error
		expectedError string
	}{
		{
			name:          "404 Not Found",
			mockErr:       fmt.Errorf("[GET /projects/{project_name}/summary][404] getProjectSummaryNotFound"),
			expectedError: "project test-project does not exist",
		},
		{
			name:          "401 Unauthorized",
			mockErr:       fmt.Errorf("[GET /projects/{project_name}/summary][401] getProjectSummaryUnauthorized"),
			expectedError: "failed to get project summary",
		},
		{
			name:          "403 Forbidden",
			mockErr:       fmt.Errorf("[GET /projects/{project_name}/summary][403] getProjectSummaryForbidden"),
			expectedError: "failed to get project summary",
		},
		{
			name:          "Connection Refused",
			mockErr:       errors.New("connection refused"),
			expectedError: "failed to get project summary: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalFunc := getProjectSummaryFunc
			defer func() { getProjectSummaryFunc = originalFunc }()

			getProjectSummaryFunc = func(projectNameOrID string, useProjectID bool) (*project.GetProjectSummaryOK, error) {
				return nil, tt.mockErr
			}

			cmd := SummaryCommand()
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs([]string{"test-project"})

			err := cmd.Execute()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}
