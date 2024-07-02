package prompt

import (
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/constants"
	aview "github.com/goharbor/harbor-cli/pkg/views/artifact/select"
	tview "github.com/goharbor/harbor-cli/pkg/views/artifact/tags/select"
	pview "github.com/goharbor/harbor-cli/pkg/views/project/select"
	rview "github.com/goharbor/harbor-cli/pkg/views/registry/select"
	repoView "github.com/goharbor/harbor-cli/pkg/views/repository/select"
	robotView "github.com/goharbor/harbor-cli/pkg/views/robot/select"
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

func GetProjectIDFromUser() int64 {
	projectID := make(chan int64)
	go func() {
		response, _ := api.ListProject()
		pview.ProjectListID(response.Payload, projectID)
	}()

	return <-projectID
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

func GetRobotPermissionsFromUser() []models.Permission {
	permissions := make(chan []models.Permission)
	go func() {
		response, _ := api.GetPermissions()
		robotView.ListPermissions(response.Payload, permissions)
	}()
	return <-permissions
}

func GetRobotIDFromUser(projectID int64) int64 {
	robotID := make(chan int64)
	var opts api.ListFlags
	opts.Q = constants.ProjectQString + strconv.FormatInt(projectID, 10)

	go func() {
		response, _ := api.ListRobot(opts)
		robotView.ListRobot(response.Payload, robotID)
	}()
	return <-robotID
}
