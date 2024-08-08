package prompt

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	aview "github.com/goharbor/harbor-cli/pkg/views/artifact/select"
	tview "github.com/goharbor/harbor-cli/pkg/views/artifact/tags/select"
	pview "github.com/goharbor/harbor-cli/pkg/views/project/select"
	rview "github.com/goharbor/harbor-cli/pkg/views/registry/select"
	repoView "github.com/goharbor/harbor-cli/pkg/views/repository/select"
	uview "github.com/goharbor/harbor-cli/pkg/views/user/select"
	iview "github.com/goharbor/harbor-cli/pkg/views/immutable/select"
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

func GetImmutableTagRule(projectName string) int64 {
	immutableid := make(chan int64)
	go func() {
		response, _ := api.ListImmutable(projectName)
		iview.ImmutableList(response.Payload, immutableid)
	}()
	return <-immutableid
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
