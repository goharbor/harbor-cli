package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	aview "github.com/goharbor/harbor-cli/pkg/views/artifact/select"
	tview "github.com/goharbor/harbor-cli/pkg/views/artifact/tags/select"
	pview "github.com/goharbor/harbor-cli/pkg/views/project/select"
	rview "github.com/goharbor/harbor-cli/pkg/views/registry/select"
	repoView "github.com/goharbor/harbor-cli/pkg/views/repository/select"
	uview "github.com/goharbor/harbor-cli/pkg/views/user/select"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Returns Harbor v2 client for given clientConfig

func PrintPayloadInJSONFormat(payload any) {
	if payload == nil {
		return
	}

	jsonStr, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonStr))
}

func GetRegistryNameFromUser() int64 {
	registryId := make(chan int64)
	go func() {
		credentialName := viper.GetString("current-credential-name")
		client := GetClientByCredentialName(credentialName)
		ctx := context.Background()
		response, err := client.Registry.ListRegistries(ctx, &registry.ListRegistriesParams{})
		if err != nil {
			log.Fatal(err)
		}

		rview.RegistryList(response.Payload, registryId)

	}()

	return <-registryId

}

func GetProjectNameFromUser() string {
	projectName := make(chan string)
	go func() {
		credentialName := viper.GetString("current-credential-name")
		client := GetClientByCredentialName(credentialName)
		ctx := context.Background()
		response, err := client.Project.ListProjects(ctx, &project.ListProjectsParams{})
		if err != nil {
			log.Fatal(err)
		}
		pview.ProjectList(response.Payload, projectName)

	}()

	return <-projectName
}

func GetRepoNameFromUser(projectName string) string {
	repositoryName := make(chan string)

	go func() {
		credentialName := viper.GetString("current-credential-name")
		client := GetClientByCredentialName(credentialName)
		ctx := context.Background()
		response, err := client.Repository.ListRepositories(ctx, &repository.ListRepositoriesParams{ProjectName: projectName})
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
		credentialName := viper.GetString("current-credential-name")
		client := GetClientByCredentialName(credentialName)
		ctx := context.Background()
		response, _ := client.Artifact.ListArtifacts(ctx, &artifact.ListArtifactsParams{ProjectName: projectName, RepositoryName: repositoryName})

		aview.ListArtifacts(response.Payload, reference)

	}()
	return <-reference
}

func ParseProjectRepo(projectRepo string) (string, string) {
	split := strings.Split(projectRepo, "/")
	if len(split) != 2 {
		log.Fatalf("invalid project/repository format: %s", projectRepo)
	}
	return split[0], split[1]
}

func ParseProjectRepoReference(projectRepoReference string) (string, string, string) {
	split := strings.Split(projectRepoReference, "/")
	if len(split) != 3 {
		log.Fatalf("invalid project/repository/reference format: %s", projectRepoReference)
	}
	return split[0], split[1], split[2]
}

func GetUserIdFromUser() int64 {
	userId := make(chan int64)

	go func() {
		credentialName := viper.GetString("current-credential-name")
		client := GetClientByCredentialName(credentialName)
		ctx := context.Background()
		response, err := client.User.ListUsers(ctx, &user.ListUsersParams{})
		if err != nil {
			log.Fatal(err)
		}
		uview.UserList(response.Payload, userId)
	}()

	return <-userId

}

func GetTagFromUser(repoName, projectName, reference string) string {
	tag := make(chan string)
	go func() {
		credentialName := viper.GetString("current-credential-name")
		client := GetClientByCredentialName(credentialName)
		ctx := context.Background()
		response, err := client.Artifact.ListTags(ctx, &artifact.ListTagsParams{ProjectName: projectName, RepositoryName: repoName, Reference: reference})
		if err != nil {
			log.Fatal(err)
		}
		tview.ListTags(response.Payload, tag)
	}()
	return <-tag
}
