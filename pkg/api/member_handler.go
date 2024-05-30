package api

import (
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
func ListMembers(projectName string) (*member.ListProjectMembersOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Member.ListProjectMembers(
		ctx,
		&member.ListProjectMembersParams{ProjectNameOrID: projectName},
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
			ProjectMember: &models.ProjectMember{
				RoleID:      int64(opts.RoleID + 1),
				MemberUser:  opts.MemberUser,
				MemberGroup: opts.MemberGroup,
			},
			ProjectNameOrID: opts.ProjectNameOrID,
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

func DeleteAllMember(projectName string) {
	var wg sync.WaitGroup
	response, _ := ListMembers(projectName)
	length := len(response.Payload)
	errChan := make(chan error, length)

	for _, member := range response.Payload {
		wg.Add(1)
		go func(memberID int64) {
			defer wg.Done()
			err := DeleteMember(projectName, memberID)
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

func DeleteMember(projectName string, memberID int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Member.DeleteProjectMember(
		ctx,
		&member.DeleteProjectMemberParams{ProjectNameOrID: projectName, Mid: memberID},
	)
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
		&member.GetProjectMemberParams{ProjectNameOrID: opts.ProjectNameOrID, Mid: opts.ID},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
