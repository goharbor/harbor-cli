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
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
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
	split := strings.SplitN(projectRepo, "/", 2) //splits only at first slash
	if len(split) != 2 {
		log.Fatalf("invalid project/repository format: %s", projectRepo)
	}
	return split[0], split[1]
}

func ParseProjectRepoReference(projectRepoReference string) (string, string, string) {
	log.Infof("Parsing input: %s", projectRepoReference)

	var ref string
	var repoPath string

	if strings.Contains(projectRepoReference, "@") {
		parts := strings.SplitN(projectRepoReference, "@", 2)
		repoPath = parts[0]
		ref = parts[1]
	} else if strings.Contains(projectRepoReference, ":") {
		lastColon := strings.LastIndex(projectRepoReference, ":")
		repoPath = projectRepoReference[:lastColon]
		ref = projectRepoReference[lastColon+1:]
	} else {
		log.Fatalf("Invalid reference format: %s", projectRepoReference)
	}

	projectRepoParts := strings.SplitN(repoPath, "/", 2)
	if len(projectRepoParts) != 2 {
		log.Fatalf("Invalid format, expected <project>/<repository>:<tag> or <project>/<repository>@<digest>, got: %s", projectRepoReference)
	}

	project := projectRepoParts[0]
	repo := projectRepoParts[1]

	return project, repo, ref
}

func SanitizeServerAddress(server string) string {
	var sb strings.Builder
	prevDash := false
	for _, r := range server {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
			prevDash = false
		} else if !prevDash {
			sb.WriteRune('-')
			prevDash = true
		}
	}

	sanitized := sb.String()
	sanitized = strings.Trim(sanitized, "-")

	return sanitized
}

func DefaultCredentialName(username, server string) string {
	sanitized := SanitizeServerAddress(server)
	return fmt.Sprintf("%s@%s", username, sanitized)
}

func SavePayloadJSON(filename string, payload any) {
	// Marshal the payload into a JSON string with indentation
	jsonStr, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}
	// Define the filename
	filename = filename + ".json"
	err = os.WriteFile(filename, jsonStr, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON data has been written to %s\n", filename)
}

// Get Password as Stdin
func GetSecretStdin(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // move to the next line after input
	return strings.TrimSpace(string(bytePassword)), nil
}
