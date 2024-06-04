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
	"regexp"
	"os"
	"strings"
	"syscall"
	"unicode"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"golang.org/x/term"
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
	split := strings.Split(projectRepoReference, "/")
	if len(split) != 3 {
		log.Fatalf("invalid project/repository/reference format: %s", projectRepoReference)
	}
	return split[0], split[1], split[2]
}

func SanitizeServerAddress(server string) string {
	re := regexp.MustCompile(`^https?://`)
	server = re.ReplaceAllString(server, "")
	re = regexp.MustCompile(`[^a-zA-Z0-9]`)
	server = re.ReplaceAllString(server, "-")
	return server
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
