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

package member

import (
	"bytes"
	"errors"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/system"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/member/create"
	"github.com/stretchr/testify/assert"
)

func TestCreateMemberCommand_Structure(t *testing.T) {
	cmd := CreateMemberCommand()

	assert.Equal(t, "create", cmd.Use)
	assert.Equal(t, "create project member", cmd.Short)
	assert.Equal(t, "create project member by Name", cmd.Long)
	assert.Equal(t, "  harbor project member create my-project --username user --role Developer", cmd.Example)
}

func TestCreateMemberCommand_Flags(t *testing.T) {
	cmd := CreateMemberCommand()
	flags := cmd.Flags()

	tests := []struct {
		name      string
		shorthand string
		defValue  string
	}{
		{"id", "", "false"},
		{"role", "", ""},
		{"roleid", "", "0"},
		{"username", "", ""},
		{"groupname", "", ""},
		{"ldapdn", "", ""},
		{"groupid", "", "0"},
		{"userid", "", "0"},
		{"grouptype", "", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.name)
			assert.NotNil(t, flag, "Flag %s should exist", tt.name)
			assert.Equal(t, tt.shorthand, flag.Shorthand)
			assert.Equal(t, tt.defValue, flag.DefValue)
		})
	}
}

func TestCreateMemberCommand_MaxArgs(t *testing.T) {
	cmd := CreateMemberCommand()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"project1", "project2"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts at most 1 arg(s)")
}

func TestCreateMemberCommand_CreationError(t *testing.T) {
	origGetProject := getProjectAPI
	origGetSysInfo := getSystemInfoAPI
	origCreateMember := createMemberAPI
	defer func() {
		getProjectAPI = origGetProject
		getSystemInfoAPI = origGetSysInfo
		createMemberAPI = origCreateMember
	}()

	authMode := "db_auth"
	getProjectAPI = func(projectNameOrID string, isID bool) (*models.Project, error) {
		return &models.Project{Name: "test-project"}, nil
	}
	getSystemInfoAPI = func() (*system.GetSystemInfoOK, error) {
		return &system.GetSystemInfoOK{
			Payload: &models.SystemInfo{AuthMode: &authMode},
		}, nil
	}
	createMemberAPI = func(opts create.CreateView) error {
		return errors.New("network timeout")
	}

	cmd := CreateMemberCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"test-project", "--username", "testuser", "--roleid", "2"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create member: network timeout")
	assert.NotContains(t, buf.String(), "successfully added user")
}

func TestCreateMemberCommand_Success(t *testing.T) {
	origGetProject := getProjectAPI
	origGetSysInfo := getSystemInfoAPI
	origCreateMember := createMemberAPI
	defer func() {
		getProjectAPI = origGetProject
		getSystemInfoAPI = origGetSysInfo
		createMemberAPI = origCreateMember
	}()

	authMode := "db_auth"
	getProjectAPI = func(projectNameOrID string, isID bool) (*models.Project, error) {
		return &models.Project{Name: "test-project"}, nil
	}
	getSystemInfoAPI = func() (*system.GetSystemInfoOK, error) {
		return &system.GetSystemInfoOK{
			Payload: &models.SystemInfo{AuthMode: &authMode},
		}, nil
	}
	var calledWith create.CreateView
	createMemberAPI = func(opts create.CreateView) error {
		calledWith = opts
		return nil
	}

	cmd := CreateMemberCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"test-project", "--username", "testuser", "--roleid", "2"})

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "test-project", calledWith.ProjectName)
	assert.Equal(t, "testuser", calledWith.MemberUser.Username)
	assert.Equal(t, 2, calledWith.RoleID)
}

func TestCreateMemberCommand_SystemInfoError(t *testing.T) {
	origGetProject := getProjectAPI
	origGetSysInfo := getSystemInfoAPI
	defer func() {
		getProjectAPI = origGetProject
		getSystemInfoAPI = origGetSysInfo
	}()

	getProjectAPI = func(projectNameOrID string, isID bool) (*models.Project, error) {
		return &models.Project{Name: "test-project"}, nil
	}
	getSystemInfoAPI = func() (*system.GetSystemInfoOK, error) {
		return nil, errors.New("unauthorized")
	}

	cmd := CreateMemberCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"test-project", "--username", "testuser", "--roleid", "2"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not access server info: unauthorized")
}
