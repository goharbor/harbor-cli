package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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

func PrintPayloadInYAMLFormat(payload any) {
	if payload == nil {
		return
	}

	yamlStr, err := yaml.Marshal(payload)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(yamlStr))
}

func ParseProjectRepo(projectRepo string) (string, string) {
	split := strings.Split(projectRepo, "/")
	if len(split) != 2 {
		log.Fatalf("invalid project/repository format: %s", projectRepo)
	}
	return split[0], split[1]
}

func ParseProjectRepoReference(projectRepoReference string) (string, string, string) {
	var projectname []string
	var reponame []string
	var referencename []string
	indexSlash := strings.Index(projectRepoReference, "/")
	indexAt := strings.Index(projectRepoReference, "@")

	if indexSlash != -1 && indexAt != -1 && indexSlash < indexAt {
		projectname = strings.Split(projectRepoReference, "/")
		reponame = strings.Split(projectname[1], "@")
		referencename = strings.Split(projectname[1], "@")
	} else {
		log.Fatalf("invalid project/repository@reference format: %s", projectRepoReference)
	}
	return projectname[0], reponame[0], referencename[1]
}

func SanitizeServerAddress(server string) string {
	re := regexp.MustCompile(`^https?://`)
	server = re.ReplaceAllString(server, "")
	re = regexp.MustCompile(`[^a-zA-Z0-9]`)
	server = re.ReplaceAllString(server, "-")
	return server
}
