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

//go:build e2e

package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	helpers "github.com/goharbor/harbor-cli/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// login runs `harbor login <server> -u <username> -p <password>` against a
// config directory scoped to this test.
func login(t *testing.T, server, username, password string) error {
	t.Helper()

	data := helpers.Initialize(t, t.TempDir())
	t.Cleanup(func() { helpers.ConfigCleanup(t, data) })

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{server})
	require.NoError(t, cmd.Flags().Set("username", username))
	require.NoError(t, cmd.Flags().Set("password", password))

	return cmd.Execute()
}

func Test_Login_Success(t *testing.T) {
	assert.NoError(t, login(t, Server(), Username(), Password()))
}

func Test_Login_Success_TrailingSlash(t *testing.T) {
	assert.NoError(t, login(t, Server()+"/", Username(), Password()))
}

func Test_Login_Success_RobotAccount(t *testing.T) {
	created := createRobot(t, fmt.Sprintf("cli-e2e-%d", time.Now().UnixNano()))

	assert.NoError(t, login(t, Server(), created.Name, created.Secret))
}

func Test_Login_Failure_WrongUsername(t *testing.T) {
	assert.Error(t, login(t, Server(), "does-not-exist", Password()))
}

func Test_Login_Failure_WrongPassword(t *testing.T) {
	assert.Error(t, login(t, Server(), Username(), "wrong"))
}
