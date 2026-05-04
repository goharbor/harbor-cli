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
package prompt

import (
	"errors"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/immutable"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/member"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/preheat"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/quota"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/scanner"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestGetRegistryNameFromUser_APIFailure(t *testing.T) {
	original := listRegistriesFunc
	defer func() { listRegistriesFunc = original }()

	listRegistriesFunc = func(opts ...api.ListFlags) (*registry.ListRegistriesOK, error) {
		return nil, errors.New("connection refused")
	}

	id, err := GetRegistryNameFromUser()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
	assert.Equal(t, int64(0), id)
}

func TestGetReferenceFromUser_APIFailure(t *testing.T) {
	original := listArtifactFunc
	defer func() { listArtifactFunc = original }()

	listArtifactFunc = func(projectName, repoName string, opts ...api.ListFlags) (artifact.ListArtifactsOK, error) {
		return artifact.ListArtifactsOK{}, errors.New("server error")
	}

	ref, err := GetReferenceFromUser("repo", "project")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server error")
	assert.Empty(t, ref)
}

func TestGetImmutableTagRule_APIFailure(t *testing.T) {
	original := listImmutableFunc
	defer func() { listImmutableFunc = original }()

	listImmutableFunc = func(projectName string) (immutable.ListImmuRulesOK, error) {
		return immutable.ListImmuRulesOK{}, errors.New("not found")
	}

	id, err := GetImmutableTagRule("test-project")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.Equal(t, int64(0), id)
}

func TestGetTagFromUser_APIFailure(t *testing.T) {
	original := listTagsFunc
	defer func() { listTagsFunc = original }()

	listTagsFunc = func(projectName, repoName, reference string) (*artifact.ListTagsOK, error) {
		return nil, errors.New("timeout")
	}

	tagName, err := GetTagFromUser("repo", "project", "ref")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
	assert.Empty(t, tagName)
}

func TestGetScannerIdFromUser_APIFailure(t *testing.T) {
	original := listScannersFunc
	defer func() { listScannersFunc = original }()

	listScannersFunc = func() (scanner.ListScannersOK, error) {
		return scanner.ListScannersOK{}, errors.New("unauthorized")
	}

	id, err := GetScannerIdFromUser()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	assert.Empty(t, id)
}

func TestGetInstanceFromUser_APIFailure(t *testing.T) {
	original := listInstanceFunc
	defer func() { listInstanceFunc = original }()

	listInstanceFunc = func(opts ...api.ListFlags) (*preheat.ListInstancesOK, error) {
		return nil, errors.New("connection refused")
	}

	name, err := GetInstanceFromUser()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
	assert.Empty(t, name)
}

func TestGetMemberIDFromUser_APIFailure(t *testing.T) {
	original := listMembersForPromptFunc
	defer func() { listMembersForPromptFunc = original }()

	listMembersForPromptFunc = func(projectNameOrID, memberName string, isName bool) (*member.ListProjectMembersOK, error) {
		return nil, errors.New("server down")
	}

	id, err := GetMemberIDFromUser("project", "member")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server down")
	assert.Equal(t, int64(0), id)
}

func TestGetMemberIDFromUser_EmptyList(t *testing.T) {
	original := listMembersForPromptFunc
	defer func() { listMembersForPromptFunc = original }()

	listMembersForPromptFunc = func(projectNameOrID, memberName string, isName bool) (*member.ListProjectMembersOK, error) {
		return &member.ListProjectMembersOK{Payload: []*models.ProjectMemberEntity{}}, nil
	}

	id, err := GetMemberIDFromUser("project", "member")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), id)
}

func TestGetQuotaIDFromUser_APIFailure(t *testing.T) {
	original := listQuotaFunc
	defer func() { listQuotaFunc = original }()

	listQuotaFunc = func(opts api.ListQuotaFlags) (*quota.ListQuotasOK, error) {
		return nil, errors.New("quota service unavailable")
	}

	id, err := GetQuotaIDFromUser()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quota service unavailable")
	assert.Equal(t, int64(0), id)
}
