package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	rview "github.com/goharbor/harbor-cli/pkg/views/registry/select"
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

func GetRegistryTypeFromUser() string {
	registryType := make(chan string)
	go func() {
		credentialName := viper.GetString("current-credential-name")
		client := GetClientByCredentialName(credentialName)
		ctx := context.Background()
		response, err := client.Registry.ListRegistryProviderTypes(
			ctx,
			&registry.ListRegistryProviderTypesParams{},
		)
		if err != nil {
			log.Fatal(err)
		}

		rview.RegistryListTypes(response.Payload, registryType)
	}()

	return <-registryType
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
