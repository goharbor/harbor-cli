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
	"fmt"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/views/project/create"
	"github.com/stretchr/testify/assert"
)

type mockProjectCreator struct {
	projectName  map[string]struct{}
	errInFilling bool
}

// FillProjectView simulates a user filling the project details through the TUI
func (m *mockProjectCreator) FillProjectView(createView *create.CreateView) error {
	if m.errInFilling {
		return fmt.Errorf("unexpected error while filling the project view")
	}
	randomProjectName := "testProject999"
	createView.ProjectName = randomProjectName
	createView.StorageLimit = "-1"
	return nil
}
func (m *mockProjectCreator) CreateProject(opts create.CreateView) error {
	if _, found := m.projectName[opts.ProjectName]; found {
		return fmt.Errorf("project with same name already exists")
	}
	m.projectName[opts.ProjectName] = struct{}{}
	return nil
}
func TestFillCreateView(t *testing.T) {
	tests := []struct {
		name         string
		input        *create.CreateView
		errInFilling bool
		expectError  bool
		expectName   string
		expectLimit  string
	}{
		{
			name:        "nil createView gets defaults then filled",
			input:       nil,
			expectName:  "testProject999",
			expectLimit: "-1",
		},
		{
			name: "empty createView gets filled",
			input: &create.CreateView{
				ProjectName:  "",
				StorageLimit: "",
			},
			expectName:  "testProject999",
			expectLimit: "-1",
		},
		{
			name:         "error in FillProjectView propagates",
			input:        &create.CreateView{},
			errInFilling: true,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectCreator{
				projectName:  make(map[string]struct{}),
				errInFilling: tt.errInFilling,
			}

			err := fillCreateView(mock, tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.input != nil {
					assert.Equal(t, tt.expectName, tt.input.ProjectName)
					assert.Equal(t, tt.expectLimit, tt.input.StorageLimit)
				}
			}
		})
	}
}
func TestCreateProject(t *testing.T) {
	projectsAreEqual := func(u1, u2 []*create.CreateView) bool {
		if len(u1) != len(u2) {
			return false
		}
		mp := make(map[string]int)
		for _, proj := range u1 {
			mp[proj.ProjectName]++
		}
		for _, proj := range u2 {
			mp[proj.ProjectName]--
		}
		for _, val := range mp {
			if val != 0 {
				return false
			}
		}
		return true
	}

	type input struct {
		opts *create.CreateView
		args []string
	}
	tests := []struct {
		name             string
		setup            func() ([]input, *mockProjectCreator)
		expectedErr      string
		expectedProjects []*create.CreateView
	}{
		{
			name: "successfully create project with all flags",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
					{
						opts: &create.CreateView{
							Public:       true,
							StorageLimit: "100",
						},
						args: []string{"my-project"},
					},
				}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedProjects: []*create.CreateView{
				{ProjectName: "my-project"},
			},
		},
		{
			name: "missing fields triggers interactive view",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
					{
						opts: &create.CreateView{},
					},
				}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedProjects: []*create.CreateView{
				{ProjectName: "testProject999"},
			},
		},
		{
			name: "project name from args with empty storage triggers interactive view",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
					{
						opts: &create.CreateView{},
						args: []string{"arg-project"},
					},
				}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedProjects: []*create.CreateView{
				{ProjectName: "testProject999"},
			},
		},
		{
			name: "proxy cache without registry id returns error",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
					{
						opts: &create.CreateView{
							StorageLimit: "100",
							ProxyCache:   true,
						},
						args: []string{"proxy-project"},
					},
				}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedErr:      "proxy cache selected but no registry ID provided",
			expectedProjects: []*create.CreateView{},
		},
		{
			name: "registry id without proxy cache returns error",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
					{
						opts: &create.CreateView{
							StorageLimit: "100",
							RegistryID:   "5",
						},
						args: []string{"bad-config"},
					},
				}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedErr:      "registry ID should only be provided when proxy-cache is enabled",
			expectedProjects: []*create.CreateView{},
		},
		{
			name: "proxy cache with registry id succeeds",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
					{
						opts: &create.CreateView{
							StorageLimit: "200",
							ProxyCache:   true,
							RegistryID:   "10",
						},
						args: []string{"cache-project"},
					},
				}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedProjects: []*create.CreateView{
				{ProjectName: "cache-project"},
			},
		},
		{
			name: "duplicate project name fails second create",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
					{opts: &create.CreateView{StorageLimit: "100"}, args: []string{"dup"}},
					{opts: &create.CreateView{StorageLimit: "100"}, args: []string{"dup"}},
				}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedErr: "failed to create project",
			expectedProjects: []*create.CreateView{
				{ProjectName: "dup"},
			},
		},
		{
			name: "create multiple projects",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
					{opts: &create.CreateView{StorageLimit: "10"}, args: []string{"proj-a"}},
					{opts: &create.CreateView{StorageLimit: "20"}, args: []string{"proj-b"}},
					{opts: &create.CreateView{StorageLimit: "30"}, args: []string{"proj-c"}},
				}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedProjects: []*create.CreateView{
				{ProjectName: "proj-a"},
				{ProjectName: "proj-b"},
				{ProjectName: "proj-c"},
			},
		},
		{
			name: "interactive view error propagates",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{
						{opts: &create.CreateView{}},
					}, &mockProjectCreator{
						projectName:  make(map[string]struct{}),
						errInFilling: true,
					}
			},
			expectedErr:      "Failed to get the required params to create project",
			expectedProjects: []*create.CreateView{},
		},
		{
			name: "no projects to create",
			setup: func() ([]input, *mockProjectCreator) {
				return []input{}, &mockProjectCreator{projectName: make(map[string]struct{})}
			},
			expectedProjects: []*create.CreateView{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputs, m := tt.setup()
			var lastErr error
			for _, in := range inputs {
				var buf bytes.Buffer
				err := CreateProject(&buf, m, in.opts, in.args)
				if err != nil {
					lastErr = err
				}
			}

			if tt.expectedErr != "" {
				assert.Error(t, lastErr)
				assert.Contains(t, lastErr.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, lastErr)
			}

			// Build the list of created projects from the mock's map
			var created []*create.CreateView
			for name := range m.projectName {
				created = append(created, &create.CreateView{ProjectName: name})
			}

			if !projectsAreEqual(created, tt.expectedProjects) {
				t.Errorf("Projects mismatch.\nExpected: %+v\nGot: %+v", tt.expectedProjects, created)
			}
		})
	}
}

func TestCreateProjectCommand(t *testing.T) {
	cmd := CreateProjectCommand()

	assert.Equal(t, "create [project name]", cmd.Use)
	assert.Equal(t, "create project", cmd.Short)
	assert.NotNil(t, cmd.Args, "Args validator should be set")
	assert.NotNil(t, cmd.RunE, "RunE should be set")

	publicFlag := cmd.Flags().Lookup("public")
	assert.NotNil(t, publicFlag)
	assert.Equal(t, "false", publicFlag.DefValue)

	registryIDFlag := cmd.Flags().Lookup("registry-id")
	assert.NotNil(t, registryIDFlag)
	assert.Equal(t, "", registryIDFlag.DefValue)

	storageLimitFlag := cmd.Flags().Lookup("storage-limit")
	assert.NotNil(t, storageLimitFlag)
	assert.Equal(t, "", storageLimitFlag.DefValue)

	proxyCacheFlag := cmd.Flags().Lookup("proxy-cache")
	assert.NotNil(t, proxyCacheFlag)
	assert.Equal(t, "false", proxyCacheFlag.DefValue)
}
