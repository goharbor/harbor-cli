package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/goharbor/go-client/pkg/harbor"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	pview "github.com/goharbor/harbor-cli/pkg/views/project/select"
	rview "github.com/goharbor/harbor-cli/pkg/views/registry/select"
	repoView "github.com/goharbor/harbor-cli/pkg/views/repository/select"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Returns Harbor v2 client for given clientConfig
func GetClientByConfig(clientConfig *harbor.ClientSetConfig) *v2client.HarborAPI {
	cs, err := harbor.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}
	return cs.V2()
}

// Returns Harbor v2 client after resolving the credential name
func GetClientByCredentialName(credentialName string) *v2client.HarborAPI {
	credential, err := GetCredentials(credentialName)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	clientConfig := &harbor.ClientSetConfig{
		URL:      credential.ServerAddress,
		Username: credential.Username,
		Password: credential.Password,
	}
	return GetClientByConfig(clientConfig)
}

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

func ParseProjectRepo(projectRepo string) (string, string) {
	split := strings.Split(projectRepo, "/")
	if len(split) != 2 {
		log.Fatalf("invalid project/repository format: %s", projectRepo)
	}
	return split[0], split[1]
}
