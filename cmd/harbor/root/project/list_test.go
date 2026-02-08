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
	"reflect"
	"runtime"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/stretchr/testify/assert"
)

// api.ListProject, api.ListAllProjects
// -> take care of public and private
// -> project store and return projects according to pagination similar to user pagination,
type MockProjectLister struct {
	projectsCnt int
	projects    []*models.Project
	expectError bool
}

func (m *MockProjectLister) mockListFunc(opts ...api.ListFlags) (project.ListProjectsOK, error) {
	res := &project.ListProjectsOK{}
	if m.expectError {
		return *res, fmt.Errorf("mock list error")
	}
	if len(opts) == 0 {
		return *res, fmt.Errorf("No options passed")
	}
	listFlags := opts[0]
	page, pageSize := listFlags.Page, listFlags.PageSize
	projects := m.populateProjects()
	lo, hi := max(pageSize*(page-1), 0), min(pageSize*page, int64(m.projectsCnt))
	if lo >= int64(m.projectsCnt) {
		return *res, nil
	}
	res.Payload = projects[lo:hi]
	return *res, nil
}

func (m *MockProjectLister) populateProjects() []*models.Project {
	projects := make([]*models.Project, 0, m.projectsCnt)
	for i := 0; i < int(m.projectsCnt); i++ {
		user := &models.Project{
			ProjectID: int32(i + 1), // #nosec G115
			Name:      fmt.Sprintf("Project%d", i+1),
		}
		projects = append(projects, user)
	}
	m.projects = projects
	return projects
}

func TestBuildListOptions(t *testing.T) {
	//basically check opts.Public, opts.Private, opts.Q and the name of the function that is returned
	getFuncName := func(i interface{}) string {
		if i == nil {
			return "nil"
		}
		return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	}
	type testInputs struct {
		private, public      bool
		opts                 *api.ListFlags
		fuzzy, match, ranges []string
	}
	tests := []struct {
		name         string
		setup        func() *testInputs
		wantedErr    string
		wantedQparam string
		wantFunc     string
		wantPublic   bool
	}{
		{
			name: "private flag returns ListProject and sets Public to false",
			setup: func() *testInputs {
				return &testInputs{
					private: true,
					opts:    &api.ListFlags{},
				}
			},
			wantFunc:   getFuncName(api.ListProject),
			wantPublic: false,
		},
		{
			name: "public flag returns ListProject and sets Public to true",
			setup: func() *testInputs {
				return &testInputs{
					public: true,
					opts:   &api.ListFlags{},
				}
			},
			wantFunc:   getFuncName(api.ListProject),
			wantPublic: true,
		},
		{
			name: "neither flag returns ListAllProjects",
			setup: func() *testInputs {
				return &testInputs{
					opts: &api.ListFlags{},
				}
			},
			wantFunc: getFuncName(api.ListAllProjects),
		},
		{
			name: "both private and public flags returns error",
			setup: func() *testInputs {
				return &testInputs{
					private: true,
					public:  true,
					opts:    &api.ListFlags{},
				}
			},
			wantedErr: "Cannot specify both --private and --public",
		},
		{
			name: "page size exceeds maximum",
			setup: func() *testInputs {
				return &testInputs{
					opts: &api.ListFlags{PageSize: 101},
				}
			},
			wantedErr: "page size should be greater than or equal to 0 and less than or equal to 100",
		},
		{
			name: "page size is negative",
			setup: func() *testInputs {
				return &testInputs{
					opts: &api.ListFlags{PageSize: -1},
				}
			},
			wantedErr: "page size should be greater than or equal to 0 and less than or equal to 100",
		},
		{
			name: "fuzzy match builds query param",
			setup: func() *testInputs {
				return &testInputs{
					opts:  &api.ListFlags{},
					fuzzy: []string{"name=test"},
				}
			},
			wantFunc:     getFuncName(api.ListAllProjects),
			wantedQparam: "name=~test",
		},
		{
			name: "exact match builds query param",
			setup: func() *testInputs {
				return &testInputs{
					opts:  &api.ListFlags{},
					match: []string{"name=myproject"},
				}
			},
			wantFunc:     getFuncName(api.ListAllProjects),
			wantedQparam: "name=myproject",
		},
		{
			name: "range builds query param",
			setup: func() *testInputs {
				return &testInputs{
					opts:   &api.ListFlags{},
					ranges: []string{"project_id=1~10"},
				}
			},
			wantFunc:     getFuncName(api.ListAllProjects),
			wantedQparam: "project_id=[1~10]",
		},
		{
			name: "multiple query params combined",
			setup: func() *testInputs {
				return &testInputs{
					opts:   &api.ListFlags{},
					fuzzy:  []string{"name=test"},
					match:  []string{"public=true"},
					ranges: []string{"project_id=1~10"},
				}
			},
			wantFunc:     getFuncName(api.ListAllProjects),
			wantedQparam: "name=~test,public=true,project_id=[1~10]",
		},
		{
			name: "invalid fuzzy key returns error",
			setup: func() *testInputs {
				return &testInputs{
					opts:  &api.ListFlags{},
					fuzzy: []string{"invalid_key=test"},
				}
			},
			wantedErr: "invalid key for query",
		},
		{
			name: "invalid fuzzy format returns error",
			setup: func() *testInputs {
				return &testInputs{
					opts:  &api.ListFlags{},
					fuzzy: []string{"badformat"},
				}
			},
			wantedErr: "invalid fuzzy arg",
		},
		{
			name: "private flag with fuzzy query",
			setup: func() *testInputs {
				return &testInputs{
					private: true,
					opts:    &api.ListFlags{},
					fuzzy:   []string{"name=test"},
				}
			},
			wantFunc:     getFuncName(api.ListProject),
			wantPublic:   false,
			wantedQparam: "name=~test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.setup()
			gotFunc, err := BuildListOptions(in.private, in.public, in.opts, in.fuzzy, in.match, in.ranges)

			// Check if we expected an error but did not get one (or vice-versa)
			if (err != nil) != (tt.wantedErr != "") {
				t.Fatalf("fetchProjects() error presence mismatch: got error %v, wantError %v", err, tt.wantedErr)
			}

			if tt.wantedErr != "" {
				assert.ErrorContains(t, err, tt.wantedErr, "Expected error to contain '%s', got '%s'", tt.wantedErr, err.Error())
			} else {
				assert.Equal(t, tt.wantPublic, in.opts.Public, "Expected opts.Public to be %t but got %t", tt.wantPublic, in.opts.Public)
				assert.Equal(t, tt.wantedQparam, in.opts.Q, "Expected query param to be %s but got %s", tt.wantedQparam, in.opts.Q)
				assert.NotNil(t, gotFunc, "Expected listFunc to be non-nil")
				assert.Equal(t, tt.wantFunc, getFuncName(gotFunc), "Expected function %s but got %s", tt.wantFunc, getFuncName(gotFunc))
			}
		})
	}
}

func TestFetchProjects(t *testing.T) {
	projectsAreEqual := func(u1, u2 []*models.Project) bool {
		if len(u1) != len(u2) {
			return false
		}
		mp := make(map[int]int)
		for _, proj := range u1 {
			mp[int(proj.ProjectID)]++
		}
		for _, proj := range u2 {
			mp[int(proj.ProjectID)]--
		}
		for _, val := range mp {
			if val != 0 {
				return false
			}
		}
		return true
	}
	tests := []struct {
		name      string
		setup     func() (api.ListFlags, *MockProjectLister)
		wantedErr string
	}{
		{
			name: "fetch all projects with page size 0 (multiple pages)",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 0}, &MockProjectLister{projectsCnt: 250}
			},
		},
		{
			name: "fetch all projects when total is exactly divisible by 100",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 0}, &MockProjectLister{projectsCnt: 200}
			},
		},
		{
			name: "fetch all projects with fewer than one page",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 0}, &MockProjectLister{projectsCnt: 50}
			},
		},
		{
			name: "fetch specific page with valid page size",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 2, PageSize: 50}, &MockProjectLister{projectsCnt: 102}
			},
		},
		{
			name: "fetch first page with page size 10",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 10}, &MockProjectLister{projectsCnt: 50}
			},
		},
		{
			name: "fetch last page with partial results",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 3, PageSize: 10}, &MockProjectLister{projectsCnt: 25}
			},
		},
		{
			name: "fetch page beyond available data returns empty",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 10, PageSize: 10}, &MockProjectLister{projectsCnt: 5}
			},
		},
		{
			name: "fetch with maximum allowed page size 100",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 100}, &MockProjectLister{projectsCnt: 150}
			},
		},
		{
			name: "fetch with zero projects in database",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 10}, &MockProjectLister{projectsCnt: 0}
			},
		},
		{
			name: "fetch all with zero projects in database",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 0}, &MockProjectLister{projectsCnt: 0}
			},
		},
		{
			name: "error during single page fetch",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 10}, &MockProjectLister{projectsCnt: 50, expectError: true}
			},
			wantedErr: "mock list error",
		},
		{
			name: "error during paginated fetch all",
			setup: func() (api.ListFlags, *MockProjectLister) {
				return api.ListFlags{Page: 1, PageSize: 0}, &MockProjectLister{projectsCnt: 50, expectError: true}
			},
			wantedErr: "mock list error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, m := tt.setup()
			allProjects, err := fetchProjects(m.mockListFunc, opts)

			// Check if we expected an error but did not get one (or vice-versa)
			if (err != nil) != (tt.wantedErr != "") {
				t.Fatalf("fetchProjects() error presence mismatch: got error %v, wantError %v", err, tt.wantedErr)
			}

			if tt.wantedErr != "" {
				assert.ErrorContains(t, err, tt.wantedErr, "Expected error to contain '%s', got '%s'", tt.wantedErr, err.Error())
			} else {
				if opts.PageSize == 0 {
					if !projectsAreEqual(allProjects, m.projects) {
						t.Errorf("Expected all of the users to be returned")
					}
				} else {
					requiredPage, requiredPageSize := opts.Page, opts.PageSize
					start := max(requiredPageSize*(requiredPage-1), 0)
					end := min(requiredPageSize*requiredPage, int64(m.projectsCnt))

					if start >= int64(m.projectsCnt) {
						if len(allProjects) != 0 {
							t.Errorf("Expected empty result for page beyond data, got %d users", len(allProjects))
						}
					} else {
						if !projectsAreEqual(allProjects, m.projects[start:end]) {
							t.Errorf("Expected different set of users")
						}
					}
				}
			}
		})
	}
}
