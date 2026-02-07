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
	"strings"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
)

// api.ListProject, api.ListAllProjects
// while mocking api-> anything having a prefix equal to opts.Name should be returned
// -> take care of public and private
// -> doesnt make much sense for name to have ranges flag (public too)
// -> fuzzy checks if the val provided is a substr of any project name
// -> sort=> name -> sorts in asc order according to name, -name=> sorts in desc order, invalid type of string passed-> gives default output
// -> project store and return projects according to pagination similar to user pagination,
// in api.ListAllProject we need to have a mix of private and public projects, whereas in api.ListProject we only need to list out one kind of projects (all pages have only private or only public)
type MockProjectLister struct {
	projectsCnt int
	// publicProjects  []*models.Project
	// privateProjects []*models.Project //ProjectID and Name are useful attributes
	projects []*models.Project
}

func (m *MockProjectLister) ListProjects(opts ...api.ListFlags) (project.ListProjectsOK, error) {
	var res project.ListProjectsOK

	return res, nil
}
func (m *MockProjectLister) ListAllProjects(opts ...api.ListFlags) (project.ListProjectsOK, error) {
	var res project.ListProjectsOK
	if len(opts) == 0 {
		return res, nil
	}
	listFlag := opts[0]
	if listFlag.Name != "" {
		res.Payload = m.filterByName(listFlag.Name)
	}
	if listFlag.Sort != "" {
		if listFlag.Sort == "name" {

		}
	}
	return res, nil
}
func (m *MockProjectLister) filterByName(name string) []*models.Project {
	var filteredProjects []*models.Project
	for _, proj := range m.projects {
		if strings.HasPrefix(proj.Name, name) {
			filteredProjects = append(filteredProjects, proj)
		}
	}
	return filteredProjects
}
func (m *MockProjectLister) filterByQuery() {

}

// odd numbered are public projects
func (m *MockProjectLister) populateProjects() {
	for i := 0; i < m.projectsCnt; i++ {
		id := i/2 + (i%2)*(m.projectsCnt-1-i) // produce pattern like 0, n-1, 1, n-2...(just to verify that sorting by name works)
		proj := &models.Project{
			ProjectID: int32(id),
			Name:      fmt.Sprintf("testProject%d", id),
		}
		if i&1 == 1 {
			proj.Metadata.Public = "true" // really dont know why this is not a boolean :(
		} else {
			proj.Metadata.Public = "false"
		}
		m.projects = append(m.projects, proj)
	}
}

func TestBuildListOptions(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { //basically check opts.Q and the kind of function that is returned

		})
	}
}

func TestFetchProjects(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
