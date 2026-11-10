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
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// execute runs the command with the given arguments, discarding its output.
// Every test here must fail during flag validation, before any prompt or
// API call can run.
func execute(t *testing.T, args ...string) error {
	t.Helper()

	cmd := ListCommand()
	cmd.SetArgs(args)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	return cmd.Execute()
}

func TestListCommandRejectsBothProjectFlags(t *testing.T) {
	err := execute(t, "--project-id", "1", "--project-name", "library")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot specify both")
}

func TestListCommandRejectsNegativeProjectID(t *testing.T) {
	err := execute(t, "--project-id=-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--project-id must not be negative")
}

func TestListCommandRejectsNegativeExecutionID(t *testing.T) {
	err := execute(t, "--execution-id=-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--execution-id must not be negative")
}
