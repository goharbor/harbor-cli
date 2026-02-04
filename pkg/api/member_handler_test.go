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

package api

import (
	"errors"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/member"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAllMember_ListMembersError(t *testing.T) {
	// Save original and restore after test
	originalListMembers := listMembersFunc
	defer func() { listMembersFunc = originalListMembers }()

	// Mock ListMembers to return an error
	listMembersFunc = func(projectNameOrID, memberName string, isName bool) (*member.ListProjectMembersOK, error) {
		return nil, errors.New("connection refused")
	}

	err := DeleteAllMember("test-project", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list members")
	assert.Contains(t, err.Error(), "connection refused")
}

func TestDeleteAllMember_NoMembers(t *testing.T) {
	originalListMembers := listMembersFunc
	defer func() { listMembersFunc = originalListMembers }()

	// Mock ListMembers to return empty list
	listMembersFunc = func(projectNameOrID, memberName string, isName bool) (*member.ListProjectMembersOK, error) {
		return &member.ListProjectMembersOK{
			Payload: []*models.ProjectMemberEntity{},
		}, nil
	}

	err := DeleteAllMember("test-project", true)

	assert.NoError(t, err)
}

func TestDeleteAllMember_Success(t *testing.T) {
	originalListMembers := listMembersFunc
	originalDeleteMember := deleteMemberFunc
	defer func() {
		listMembersFunc = originalListMembers
		deleteMemberFunc = originalDeleteMember
	}()

	// Mock ListMembers to return two members
	listMembersFunc = func(projectNameOrID, memberName string, isName bool) (*member.ListProjectMembersOK, error) {
		return &member.ListProjectMembersOK{
			Payload: []*models.ProjectMemberEntity{
				{ID: 1, EntityName: "user1"},
				{ID: 2, EntityName: "user2"},
			},
		}, nil
	}

	// Mock DeleteMember to succeed
	deleteMemberFunc = func(projectName string, memberID int64, xIsResourceName bool) error {
		return nil
	}

	err := DeleteAllMember("test-project", true)

	assert.NoError(t, err)
}

func TestDeleteAllMember_DeleteError(t *testing.T) {
	originalListMembers := listMembersFunc
	originalDeleteMember := deleteMemberFunc
	defer func() {
		listMembersFunc = originalListMembers
		deleteMemberFunc = originalDeleteMember
	}()

	// Mock ListMembers to return one member
	listMembersFunc = func(projectNameOrID, memberName string, isName bool) (*member.ListProjectMembersOK, error) {
		return &member.ListProjectMembersOK{
			Payload: []*models.ProjectMemberEntity{
				{ID: 1, EntityName: "user1"},
			},
		}, nil
	}

	// Mock DeleteMember to fail
	deleteMemberFunc = func(projectName string, memberID int64, xIsResourceName bool) error {
		return errors.New("unauthorized")
	}

	err := DeleteAllMember("test-project", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete")
}
