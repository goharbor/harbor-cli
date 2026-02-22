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
	"fmt"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/views/project/create"
	"github.com/stretchr/testify/assert"
)

type mockProjectCreator struct {
	projectName  map[string]struct{}
	errInFilling bool
}

func (m *mockProjectCreator) fillProjectView(createView *create.CreateView) error {
	if m.errInFilling {
		return fmt.Errorf("unexpected error while filling the project view")
	}
	randomProjectName := "testProject999"
	createView.ProjectName = randomProjectName
	createView.StorageLimit = "-1"
	return nil
}

func (m *mockProjectCreator) createProjectAPI(opts create.CreateView) error {
	if _, found := m.projectName[opts.ProjectName]; found {
		return fmt.Errorf("project with same name already exists")
	}
	m.projectName[opts.ProjectName] = struct{}{}
	return nil
}
func TestCreateProject(t *testing.T) {
	origFillProjectView := fillProjectView
	defer func() { fillProjectView = origFillProjectView }()

	projectsAreEqual := func(p1, p2 []*create.CreateView) bool {
		if len(p1) != len(p2) {
			return false
		}
		mp := make(map[string]int)
		for _, proj := range p1 {
			mp[proj.ProjectName]++
		}
		for _, proj := range p2 {
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
		inputs           []input
		errInFilling     bool
		expectedErr      string
		expectedProjects []*create.CreateView
	}{
		{
			name: "project creation with user provided project name and storage should succeed",
			inputs: []input{
				{
					opts: &create.CreateView{
						Public:       true,
						StorageLimit: "100",
					},
					args: []string{"my-project"},
				},
			},
			expectedProjects: []*create.CreateView{
				{ProjectName: "my-project"},
			},
		},
		{
			name: "project creation should succeed if no error in interactive view",
			inputs: []input{
				{opts: &create.CreateView{}},
			},
			expectedProjects: []*create.CreateView{
				{ProjectName: "testProject999"},
			},
		},
		{
			name: "duplicate project name fails second create",
			inputs: []input{
				{opts: &create.CreateView{StorageLimit: "100"}, args: []string{"dup"}},
				{opts: &create.CreateView{StorageLimit: "100"}, args: []string{"dup"}},
			},
			expectedErr: "failed to create project",
			expectedProjects: []*create.CreateView{
				{ProjectName: "dup"},
			},
		},
		{
			name: "create multiple projects",
			inputs: []input{
				{opts: &create.CreateView{StorageLimit: "10"}, args: []string{"proj-a"}},
				{opts: &create.CreateView{StorageLimit: "20"}, args: []string{"proj-b"}},
				{opts: &create.CreateView{StorageLimit: "30"}, args: []string{"proj-c"}},
			},
			expectedProjects: []*create.CreateView{
				{ProjectName: "proj-a"},
				{ProjectName: "proj-b"},
				{ProjectName: "proj-c"},
			},
		},
		{
			name: "buildCreateView error propagates without calling API",
			inputs: []input{
				{opts: &create.CreateView{}},
			},
			errInFilling:     true,
			expectedErr:      "failed to get the required params to create project",
			expectedProjects: []*create.CreateView{},
		},
		{
			name:             "no projects to create",
			inputs:           []input{},
			expectedProjects: []*create.CreateView{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockProjectCreator{
				projectName:  make(map[string]struct{}),
				errInFilling: tt.errInFilling,
			}
			fillProjectView = m.fillProjectView

			var lastErr error
			for _, in := range tt.inputs {
				err := createProject(m.createProjectAPI, in.opts, in.args)
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

func TestBuildCreateView(t *testing.T) {
	origFillProjectView := fillProjectView
	defer func() { fillProjectView = origFillProjectView }()

	tests := []struct {
		name         string
		opts         *create.CreateView
		args         []string
		errInFilling bool
		expectedErr  string
		expectedView *create.CreateView
	}{
		{
			name: "all flags provided returns opts directly",
			opts: &create.CreateView{
				ProjectName:  "full-project",
				Public:       true,
				StorageLimit: "500",
			},
			args: nil,
			expectedView: &create.CreateView{
				ProjectName:  "full-project",
				Public:       true,
				StorageLimit: "500",
			},
		},
		{
			name: "missing storage triggers interactive view",
			opts: &create.CreateView{},
			args: []string{"some-project"},
			expectedView: &create.CreateView{
				ProjectName:  "testProject999",
				StorageLimit: "-1",
			},
		},
		{
			name: "missing project name triggers interactive view",
			opts: &create.CreateView{
				StorageLimit: "100",
			},
			expectedView: &create.CreateView{
				ProjectName:  "testProject999",
				StorageLimit: "-1",
			},
		},
		{
			name: "proxy cache without registry id returns error",
			opts: &create.CreateView{
				StorageLimit: "100",
				ProxyCache:   true,
			},
			args:        []string{"proxy-proj"},
			expectedErr: "proxy cache selected but no registry ID provided",
		},
		{
			name: "registry id without proxy cache returns error",
			opts: &create.CreateView{
				StorageLimit: "100",
				RegistryID:   "5",
			},
			args:        []string{"bad-proj"},
			expectedErr: "registry ID should only be provided when proxy-cache is enabled",
		},
		{
			name: "proxy cache with registry id succeeds",
			opts: &create.CreateView{
				ProjectName:  "cache-proj",
				StorageLimit: "200",
				ProxyCache:   true,
				RegistryID:   "10",
			},
			expectedView: &create.CreateView{
				ProjectName:  "cache-proj",
				StorageLimit: "200",
				ProxyCache:   true,
				RegistryID:   "10",
			},
		},
		{
			name:         "interactive view error propagates",
			opts:         &create.CreateView{},
			errInFilling: true,
			expectedErr:  "failed to get the required params to create project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockProjectCreator{
				projectName:  make(map[string]struct{}),
				errInFilling: tt.errInFilling,
			}
			fillProjectView = m.fillProjectView

			view, err := buildCreateView(tt.opts, tt.args)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, view)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, view)
				assert.Equal(t, tt.expectedView.ProjectName, view.ProjectName)
				assert.Equal(t, tt.expectedView.StorageLimit, view.StorageLimit)
				assert.Equal(t, tt.expectedView.ProxyCache, view.ProxyCache)
				assert.Equal(t, tt.expectedView.RegistryID, view.RegistryID)
			}
		})
	}
}
