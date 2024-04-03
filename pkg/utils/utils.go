package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/goharbor/go-client/pkg/harbor"
	view "github.com/goharbor/harbor-cli/pkg/views/project/select"

	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	log "github.com/sirupsen/logrus"
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
	credential, err := resolveCredential(credentialName)
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

func GetProjectNameFromUser(credentialName string) string {
	projectName := make(chan string)
	go func() {
		client := GetClientByCredentialName(credentialName)
		ctx := context.Background()
		response, err := client.Project.ListProjects(ctx, &project.ListProjectsParams{})
		if err != nil {
			log.Fatal(err)
		}
		view.ProjectList(response.Payload, projectName)

	}()

	return <-projectName
}
