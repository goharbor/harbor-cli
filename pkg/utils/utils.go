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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
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

func ParseProjectRepoReference(projectRepoReference string) (string, string, string, error) {
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
		return "", "", "", fmt.Errorf("Invalid reference format: %s", projectRepoReference)
	}

	projectRepoParts := strings.SplitN(repoPath, "/", 2)
	if len(projectRepoParts) != 2 {
		return "", "", "", fmt.Errorf("Invalid format, expected <project>/<repository>:<tag> or <project>/<repository>@<digest>, got: %s", projectRepoReference)
	}

	project := projectRepoParts[0]
	repo := projectRepoParts[1]

	return project, repo, ref, nil
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
