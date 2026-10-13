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
package artifact_test

import (
	"bytes"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root/artifact"
	helpers "github.com/goharbor/harbor-cli/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestDeleteArtifactCommand_ProjectSelectionFailure(t *testing.T) {
	helpers.SetMockKeyring(t)
	data := helpers.Initialize(t, t.TempDir())
	defer helpers.ConfigCleanup(t, data)

	cmd := artifact.DeleteArtifactCommand()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get project name")
}
