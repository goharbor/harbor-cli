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
package retention

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteRetentionRuleCommand(t *testing.T) {
	cmd := DeleteRetentionRuleCommand()

	assert.Equal(t, "delete", cmd.Use)
	assert.Equal(t, "Delete a tag retention rule for a project", cmd.Short)
	assert.NotNil(t, cmd.RunE)
}

func TestDeleteRetentionRuleCommand_Errors(t *testing.T) {
	// Backup and restore package-level variables
	orgGetProjectName := getProjectNameFromUser
	orgGetRetentionId := getRetentionId
	orgGetRetentionTagRule := getRetentionTagRule
	orgDeleteRetention := deleteRetention

	t.Cleanup(func() {
		getProjectNameFromUser = orgGetProjectName
		getRetentionId = orgGetRetentionId
		getRetentionTagRule = orgGetRetentionTagRule
		deleteRetention = orgDeleteRetention
	})

	t.Run("No retention policy exists - success no-op", func(t *testing.T) {
		getRetentionId = func(projectNameorID string, isName bool) (string, error) {
			return "", errors.New("No retention policy exists for this project")
		}

		cmd := DeleteRetentionRuleCommand()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--project-name", "test-project"})

		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("GetRetentionId error - returns error", func(t *testing.T) {
		getRetentionId = func(projectNameorID string, isName bool) (string, error) {
			return "", errors.New("some network error")
		}

		cmd := DeleteRetentionRuleCommand()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--project-name", "test-project"})

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error retrieving retention policy ID: some network error")
	})

	t.Run("GetRetentionTagRule error - returns error", func(t *testing.T) {
		getRetentionId = func(projectNameorID string, isName bool) (string, error) {
			return "42", nil
		}
		getRetentionTagRule = func(retentionID string) (int64, error) {
			return 0, errors.New("selection aborted")
		}

		cmd := DeleteRetentionRuleCommand()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--project-name", "test-project"})

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Equal(t, "selection aborted", err.Error())
	})

	t.Run("No retention rules found - success no-op", func(t *testing.T) {
		getRetentionId = func(projectNameorID string, isName bool) (string, error) {
			return "42", nil
		}
		getRetentionTagRule = func(retentionID string) (int64, error) {
			return -1, nil
		}

		cmd := DeleteRetentionRuleCommand()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--project-name", "test-project"})

		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("DeleteRetention error - returns error", func(t *testing.T) {
		getRetentionId = func(projectNameorID string, isName bool) (string, error) {
			return "42", nil
		}
		getRetentionTagRule = func(retentionID string) (int64, error) {
			return 0, nil
		}
		deleteRetention = func(retentionID string, ruleIndex int) error {
			return errors.New("delete failed")
		}

		cmd := DeleteRetentionRuleCommand()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--project-name", "test-project"})

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})

	t.Run("Delete retention rule - success", func(t *testing.T) {
		getRetentionId = func(projectNameorID string, isName bool) (string, error) {
			return "42", nil
		}
		getRetentionTagRule = func(retentionID string) (int64, error) {
			return 0, nil
		}
		var calledRetentionID string
		var calledIndex int
		deleteRetention = func(retentionID string, ruleIndex int) error {
			calledRetentionID = retentionID
			calledIndex = ruleIndex
			return nil
		}

		cmd := DeleteRetentionRuleCommand()
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--project-name", "test-project"})

		err := cmd.Execute()
		assert.NoError(t, err)
		assert.Equal(t, "42", calledRetentionID)
		assert.Equal(t, 0, calledIndex)
	})
}
