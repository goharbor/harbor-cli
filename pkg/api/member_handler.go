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
	"fmt"
	"sync"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/member"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/member/create"
	log "github.com/sirupsen/logrus"
)

// View a Member in a project
func ListMember(opts ListMemberOptions) (*member.ListProjectMembersOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Member.ListProjectMembers(
		ctx,
		&member.ListProjectMembersParams{
			XIsResourceName: &opts.XIsResourceName,
			ProjectNameOrID: opts.ProjectNameOrID,
			Entityname:      &opts.EntityName,
			Page:            &opts.Page,
			PageSize:        &opts.PageSize,
		},
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// List Members in project
func ListMembers(projectNameOrID, memberName string, isName bool) (*member.ListProjectMembersOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Member.ListProjectMembers(
		ctx,
		&member.ListProjectMembersParams{ProjectNameOrID: projectNameOrID, XIsResourceName: &isName, Entityname: &memberName},
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Used to create a Project Member
func CreateMember(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	response, err := client.Member.CreateProjectMember(
		ctx, &member.CreateProjectMemberParams{
			XIsResourceName: &opts.XIsResourceID,
			ProjectMember: &models.ProjectMember{
				RoleID:      int64(opts.RoleID + 1),
				MemberUser:  opts.MemberUser,
				MemberGroup: opts.MemberGroup,
			},
			ProjectNameOrID: opts.ProjectName,
		},
	)
	if err != nil {
		return err
	}

	if response != nil {
		log.Info("Member created successfully")
	}

	return nil
}

func DeleteAllMember(projectName string, xIsResourceName bool) {
	var wg sync.WaitGroup
	response, _ := ListMembers(projectName, "", true)
	length := len(response.Payload)
	errChan := make(chan error, length)

	if length < 1 {
		log.Info("No members found in project")
		return
	}

	for _, member := range response.Payload {
		wg.Add(1)
		go func(memberID int64) {
			defer wg.Done()
			err := DeleteMember(projectName, memberID, xIsResourceName)
			if err != nil {
				errChan <- err
			}
		}(member.ID) // Pass member.ID to the goroutine
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Handle errors after all deletions are done
	for err := range errChan {
		if err != nil {
			log.Errorln("Error:", err)
		}
	}
}

func DeleteMember(projectName string, memberID int64, xIsResourceName bool) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Member.DeleteProjectMember(
		ctx,
		&member.DeleteProjectMemberParams{ProjectNameOrID: projectName, Mid: memberID, XIsResourceName: &xIsResourceName},
	)
	if err != nil {
		return err
	}

	log.Info("Member deleted successfully")
	return nil
}

func DeleteMemberByUsername(projectName string, username string, xIsResourceName bool) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	members, err := ListMembers(projectName, username, true)
	if err != nil {
		return err
	}

	var memberID int64
	for _, m := range members.Payload {
		if m.EntityName == username {
			memberID = m.ID
			break
		}
	}

	if memberID == 0 {
		return fmt.Errorf("member with username '%s' not found in project '%s'", username, projectName)
	}

	_, err = client.Member.DeleteProjectMember(ctx, &member.DeleteProjectMemberParams{ProjectNameOrID: projectName, Mid: memberID, XIsResourceName: &xIsResourceName})
	if err != nil {
		return err
	}

	log.Info("Member deleted successfully")
	return nil
}

func UpdateMember(opts UpdateMemberOptions) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Member.UpdateProjectMember(
		ctx,
		&member.UpdateProjectMemberParams{
			XIsResourceName: &opts.XIsResourceName,
			ProjectNameOrID: opts.ProjectNameOrID,
			Mid:             opts.ID,
			Role:            opts.RoleID,
		},
	)
	if err != nil {
		return err
	}

	log.Info("member role updated successfully")

	return nil
}

func GetMember(opts GetMemberOptions) (*member.GetProjectMemberOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Member.GetProjectMember(
		ctx,
		&member.GetProjectMemberParams{ProjectNameOrID: opts.ProjectNameOrID, Mid: opts.ID, XIsResourceName: &opts.XIsResourceName},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
