package prompt

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	aview "github.com/goharbor/harbor-cli/pkg/views/artifact/select"
	tview "github.com/goharbor/harbor-cli/pkg/views/artifact/tags/select"
	mview "github.com/goharbor/harbor-cli/pkg/views/member/select"
	pview "github.com/goharbor/harbor-cli/pkg/views/project/select"
	rview "github.com/goharbor/harbor-cli/pkg/views/registry/select"
	repoView "github.com/goharbor/harbor-cli/pkg/views/repository/select"
	uview "github.com/goharbor/harbor-cli/pkg/views/user/select"
	log "github.com/sirupsen/logrus"
)

func GetRegistryNameFromUser() int64 {
	registryId := make(chan int64)
	go func() {
		response, _ := api.ListRegistries()
		rview.RegistryList(response.Payload, registryId)
	}()

	return <-registryId
}

func GetProjectNameFromUser() string {
	projectName := make(chan string)
	go func() {
		response, _ := api.ListProject()
		pview.ProjectList(response.Payload, projectName)
	}()

	return <-projectName
}

func GetRepoNameFromUser(projectName string) string {
	repositoryName := make(chan string)

	go func() {
		response, err := api.ListRepository(projectName)
		if err != nil {
			log.Fatal(err)
		}
		repoView.RepositoryList(response.Payload, repositoryName)
	}()

	return <-repositoryName
}

// complete the function
func GetReferenceFromUser(repositoryName string, projectName string) string {
	reference := make(chan string)
	go func() {
		response, _ := api.ListArtifact(projectName, repositoryName)
		aview.ListArtifacts(response.Payload, reference)
	}()
	return <-reference
}

func GetUserIdFromUser() int64 {
	userId := make(chan int64)

	go func() {
		response, _ := api.ListUsers()
		uview.UserList(response.Payload, userId)
	}()

	return <-userId
}

func GetTagFromUser(repoName, projectName, reference string) string {
	tag := make(chan string)
	go func() {
		response, _ := api.ListTags(projectName, repoName, reference)
		tview.ListTags(response.Payload, tag)
	}()
	return <-tag
}

func GetTagNameFromUser() string {
	repoName := make(chan string)

	go func() {
	}()

	return <-repoName
}

// Get GetMemberIDFromUser choosing from list of members
func GetMemberIDFromUser(projectName string) int64 {
	memberId := make(chan int64)
	go func() {
		response, _ := api.ListMembers(projectName)
		mview.MemberList(response.Payload, memberId)
	}()

	return <-memberId
}

// Get Member Role ID selection from user
func GetRoleIDFromUser() int64 {
	roleID := make(chan int64)
	go func() {
		roles := []string{"Project Admin", "Developer", "Guest", "Maintainer", "Limited Guest"}
		mview.RoleList(roles, roleID)
	}()

	return <-roleID
}
