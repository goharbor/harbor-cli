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
package view

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var detailsColumns = []table.Column{
	{Title: "Attribute", Width: tablelist.WidthL * 2},
	{Title: "Value", Width: tablelist.Width3XL * 2},
}

func ViewArtifact(artifact *models.Artifact) {
	displayDetailTable(artifact)
}

func displayDetailTable(artifact *models.Artifact) {
	var basicRows []table.Row
	var otherRows []table.Row

	// BASIC INFO
	addRow(&basicRows, "ID", strconv.FormatInt(int64(artifact.ID), 10))
	addRow(&basicRows, "Repository ID", strconv.FormatInt(int64(artifact.RepositoryID), 10))
	addRow(&basicRows, "Digest", artifact.Digest)
	addRow(&basicRows, "Media Type", artifact.MediaType)
	addRow(&basicRows, "Type", artifact.Type)

	// Tags
	if len(artifact.Tags) > 0 {
		var tagNames []string
		for _, tag := range artifact.Tags {
			tagNames = append(tagNames, tag.Name)
		}
		addMultiLineData(&basicRows, "Tags", tagNames, ", ", 80)
	}

	addRow(&basicRows, "Size", utils.FormatSize(artifact.Size))
	pushTime, _ := utils.FormatCreatedTime(artifact.PushTime.String())
	addRow(&basicRows, "Push Time", pushTime)

	// Platform Info from ExtraAttrs
	if artifact.ExtraAttrs != nil {
		if arch, ok := artifact.ExtraAttrs["architecture"].(string); ok {
			addRow(&basicRows, "Architecture", arch)
		}
		if osVal, ok := artifact.ExtraAttrs["os"].(string); ok {
			addRow(&basicRows, "OS", osVal)
		}
		if created, ok := artifact.ExtraAttrs["created"].(string); ok {
			addRow(&basicRows, "Created", created)
		}
	}

	// Author from config
	if artifact.ExtraAttrs != nil {
		if config, ok := artifact.ExtraAttrs["config"].(map[string]any); ok {
			if author, ok := config["author"].(string); ok {
				addRow(&basicRows, "Author", author)
			}
		}
	}

	// Version (if present in ExtraAttrs)
	if artifact.ExtraAttrs != nil {
		if version, ok := artifact.ExtraAttrs["version"].(string); ok {
			addRow(&basicRows, "Version", version)
		}
		if apiVersion, ok := artifact.ExtraAttrs["apiVersion"].(string); ok {
			addRow(&basicRows, "API Version", apiVersion)
		}
		if appVersion, ok := artifact.ExtraAttrs["appVersion"].(string); ok {
			addRow(&basicRows, "App Version", appVersion)
		}
		if description, ok := artifact.ExtraAttrs["description"].(string); ok {
			addRow(&basicRows, "Description", description)
		}
	}

	// OTHER INFO
	// Config details
	if artifact.ExtraAttrs != nil {
		if config, ok := artifact.ExtraAttrs["config"].(map[string]any); ok {
			// Environment Variables
			if env, ok := config["Env"].([]any); ok && len(env) > 0 {
				var envStrings []string
				for _, e := range env {
					envStrings = append(envStrings, fmt.Sprintf("%v", e))
				}
				addMultiLineData(&otherRows, "Environment Variables", envStrings, "", 0)
			}
			// Exposed Ports
			if ports, ok := config["ExposedPorts"].(map[string]any); ok && len(ports) > 0 {
				var portsList []string
				for port := range ports {
					portsList = append(portsList, port)
				}
				addMultiLineData(&otherRows, "Exposed Ports", portsList, "", 0)
			}
			// Volumes
			if volumes, ok := config["Volumes"].(map[string]any); ok && len(volumes) > 0 {
				var volumesList []string
				for volume := range volumes {
					volumesList = append(volumesList, volume)
				}
				addMultiLineData(&otherRows, "Volumes", volumesList, "", 0)
			}
			// Entrypoint
			if entrypoint, ok := config["Entrypoint"].([]any); ok && len(entrypoint) > 0 {
				var entryList []string
				for _, e := range entrypoint {
					entryList = append(entryList, fmt.Sprintf("%v", e))
				}
				addRow(&otherRows, "Entrypoint", strings.Join(entryList, " "))
			}
			// Command
			if cmd, ok := config["Cmd"].([]any); ok && len(cmd) > 0 {
				var cmdList []string
				for _, c := range cmd {
					cmdList = append(cmdList, fmt.Sprintf("%v", c))
				}
				addRow(&otherRows, "Command", strings.Join(cmdList, " "))
			}
			// Labels
			if labels, ok := config["Labels"].(map[string]any); ok && len(labels) > 0 {
				var labelsList []string
				for k, v := range labels {
					labelsList = append(labelsList, fmt.Sprintf("%s: %v", k, v))
				}
				addMultiLineData(&otherRows, "Labels", labelsList, "", 0)
			}
		}
	}

	// Other ExtraAttrs (excluding those already shown)
	if artifact.ExtraAttrs != nil {
		for key, value := range artifact.ExtraAttrs {
			if key == "architecture" || key == "os" || key == "config" || key == "created" || key == "layers" || key == "version" || key == "description" || key == "apiVersion" || key == "appVersion" {
				continue
			}
			switch v := value.(type) {
			case string, bool, int, int64, float64:
				addRow(&otherRows, key, fmt.Sprintf("%v", v))
			default:
				jsonBytes, err := json.MarshalIndent(value, "", "  ")
				if err == nil {
					lines := strings.Split(string(jsonBytes), "\n")
					addMultiLineData(&otherRows, key, lines, "", 0)
				} else {
					addRow(&otherRows, key, "(complex structure - see JSON output)")
				}
			}
		}
	}

	if len(artifact.References) > 0 {
		for i, ref := range artifact.References {
			if ref.Platform != nil && ref.Platform.Architecture != "" {
				addRow(&otherRows, fmt.Sprintf("Reference %d", i+1), fmt.Sprintf("%s: %s", ref.Platform.Architecture, ref.ChildDigest))
			}
		}
	}

	// Additional links
	if len(artifact.AdditionLinks) > 0 {
		var links []string
		for key := range artifact.AdditionLinks {
			links = append(links, key)
		}
		addRow(&otherRows, "Available Endpoints", strings.Join(links, ", "))
	}

	// Display the tables
	fmt.Println("\nBASIC INFORMATION")
	basicTable := tablelist.NewModel(detailsColumns, basicRows, len(basicRows))
	if _, err := tea.NewProgram(basicTable).Run(); err != nil {
		fmt.Println("Error displaying basic artifact details:", err)
		os.Exit(1)
	}

	fmt.Println("\nOTHER INFORMATION")
	otherTable := tablelist.NewModel(detailsColumns, otherRows, len(otherRows))
	if _, err := tea.NewProgram(otherTable).Run(); err != nil {
		fmt.Println("Error displaying other artifact details:", err)
		os.Exit(1)
	}
}

func addRow(rows *[]table.Row, key string, value string) {
	*rows = append(*rows, table.Row{key, value})
}

func addMultiLineData(rows *[]table.Row, key string, values []string, separator string, maxWidth int) {
	if len(values) == 0 {
		return
	}

	firstValue := values[0]
	if separator != "" && maxWidth > 0 {
		if len(firstValue) > maxWidth {
			addRow(rows, key, firstValue)
		} else {
			currentLine := firstValue
			nextIndex := 1

			for nextIndex < len(values) {
				nextValue := values[nextIndex]
				if len(currentLine)+len(separator)+len(nextValue) <= maxWidth {
					currentLine += separator + nextValue
					nextIndex++
				} else {
					break
				}
			}

			addRow(rows, key, currentLine)

			currentLine = ""
			for i := nextIndex; i < len(values); i++ {
				nextValue := values[i]
				if currentLine == "" {
					currentLine = nextValue
				} else if len(currentLine)+len(separator)+len(nextValue) <= maxWidth {
					currentLine += separator + nextValue
				} else {
					addRow(rows, "", currentLine)
					currentLine = nextValue
				}
			}

			if currentLine != "" {
				addRow(rows, "", currentLine)
			}

			return
		}
	} else {
		addRow(rows, key, firstValue)
	}

	for i := 1; i < len(values); i++ {
		addRow(rows, "", values[i])
	}
}
