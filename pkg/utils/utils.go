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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/charmbracelet/bubbles/table"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	uview "github.com/goharbor/harbor-cli/pkg/views/user/select"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v4"
	"golang.org/x/term"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func ParseProjectRepo(projectRepo string) (project, repo string, err error) {
	split := strings.SplitN(projectRepo, "/", 2) // splits only at first slash
	if len(split) != 2 {
		return "", "", fmt.Errorf("invalid project/repository format: %s", projectRepo)
	}
	return split[0], split[1], nil
}

func ParseProjectRepoReference(projectRepoReference string) (project, repo, reference string, err error) {
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
		return "", "", "", fmt.Errorf("invalid reference format: %s", projectRepoReference)
	}

	projectRepoParts := strings.SplitN(repoPath, "/", 2)
	if len(projectRepoParts) != 2 {
		return "", "", "", fmt.Errorf("invalid format, expected <project>/<repository>:<tag> or <project>/<repository>@<digest>, got: %s", projectRepoReference)
	}

	project = projectRepoParts[0]
	repo = projectRepoParts[1]

	return project, repo, ref, err
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

func StorageStringToBytes(storage string) (int64, error) {
	// Define the conversion multipliers
	multipliers := map[string]int64{
		"MiB": 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"TiB": 1024 * 1024 * 1024 * 1024,
	}

	// Define the regex to parse the input string
	re := regexp.MustCompile(`^(\d+)(MiB|GiB|TiB)$`)
	matches := re.FindStringSubmatch(storage)
	if matches == nil {
		return 0, errors.New("invalid storage format")
	}

	// Extract the value and unit from the matches
	valueStr, unit := matches[1], matches[2]
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Calculate the value in bytes
	bytes := value * multipliers[unit]

	// Check if the value exceeds 1024 TB
	maxBytes := 1024 * 1024 * 1024 * 1024 * 1024
	if bytes > int64(maxBytes) {
		return 0, errors.New("value exceeds 1024 TB")
	}

	return bytes, nil
}

func SavePayloadJSON(filename string, payload any) {
	// Marshal the payload into a JSON string with indentation
	jsonStr, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}
	// Define the filename
	filename = filename + ".json"
	err = os.WriteFile(filename, jsonStr, 0o600)
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

func ToKebabCase(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "-")
}

func FromKebabCase(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		words[i] = cases.Title(language.English).String(word)
	}
	return strings.Join(words, " ")
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	return s
	// trings.ToUpper(s[:1]) + s[1:]
}

// GetUserIdFromUser retrieves the user ID from the current user context using viper and the Harbor client.
func GetUserIdFromUser() int64 {
	credentialName := viper.GetString("current-credential-name")
	client, err := GetClientByCredentialName(credentialName)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	response, err := client.User.ListUsers(ctx, &user.ListUsersParams{})
	if err != nil {
		log.Fatal(err)
	}
	userId, err := uview.UserList(response.Payload)
	if err != nil {
		log.Fatal(err)
	}
	return userId
}

// RemoveColumns removes columns with specified titles from the given columns array.
func RemoveColumns(columns []table.Column, colsToRemove []string) []table.Column {
	titleMap := make(map[string]bool)
	for _, title := range colsToRemove {
		titleMap[title] = true
	}

	var filteredColumns []table.Column
	for _, column := range columns {
		if !titleMap[column.Title] {
			filteredColumns = append(filteredColumns, column)
		}
	}

	return filteredColumns
}
