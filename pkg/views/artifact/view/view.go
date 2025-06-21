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
	var rows []table.Row

	// SECTION: Basic metadata
	addSectionHeader(&rows, "BASIC INFO")
	addRow(&rows, "ID", strconv.FormatInt(int64(artifact.ID), 10))
	addRow(&rows, "Repository ID", strconv.FormatInt(int64(artifact.RepositoryID), 10))
	addRow(&rows, "Digest", artifact.Digest)
	addRow(&rows, "Media Type", artifact.MediaType)
	addRow(&rows, "Type", artifact.Type)

	// SECTION: Tags
	if len(artifact.Tags) > 0 {
		addSectionHeader(&rows, "TAGS")
		var tagNames []string
		for _, tag := range artifact.Tags {
			tagNames = append(tagNames, tag.Name)
		}

		// Handle tags as multi-line data with actual multiple rows
		addMultiLineData(&rows, "Tags", tagNames, ", ", 80)
	}

	// SECTION: Size and Timing
	addSectionHeader(&rows, "SIZE & TIMING")
	addRow(&rows, "Size", utils.FormatSize(artifact.Size))
	pushTime, _ := utils.FormatCreatedTime(artifact.PushTime.String())
	addRow(&rows, "Push Time", pushTime)

	// SECTION: Platform Info from ExtraAttrs
	addSectionHeader(&rows, "PLATFORM")
	if artifact.ExtraAttrs != nil {
		// Architecture
		if arch, ok := artifact.ExtraAttrs["architecture"].(string); ok {
			addRow(&rows, "Architecture", arch)
		}

		// OS
		if osVal, ok := artifact.ExtraAttrs["os"].(string); ok {
			addRow(&rows, "OS", osVal)
		}

		// Created timestamp
		if created, ok := artifact.ExtraAttrs["created"].(string); ok {
			addRow(&rows, "Created", created)
		}

		// Layers
		if layers, ok := artifact.ExtraAttrs["layers"].([]any); ok {
			addRow(&rows, "Layer Count", strconv.Itoa(len(layers)))
		}
	}

	// SECTION: Config Information
	if artifact.ExtraAttrs != nil {
		if config, ok := artifact.ExtraAttrs["config"].(map[string]any); ok {
			addSectionHeader(&rows, "CONFIGURATION")

			// Author
			if author, ok := config["author"].(string); ok {
				addRow(&rows, "Author", author)
			}

			// Environment Variables
			if env, ok := config["Env"].([]any); ok && len(env) > 0 {
				// Add environment variables as separate rows
				var envStrings []string
				for _, e := range env {
					envStrings = append(envStrings, fmt.Sprintf("%v", e))
				}
				addMultiLineData(&rows, "Environment Variables", envStrings, "", 0)
			}

			// Exposed Ports
			if ports, ok := config["ExposedPorts"].(map[string]any); ok && len(ports) > 0 {
				var portsList []string
				for port := range ports {
					portsList = append(portsList, port)
				}
				addMultiLineData(&rows, "Exposed Ports", portsList, "", 0)
			}

			// Volumes
			if volumes, ok := config["Volumes"].(map[string]any); ok && len(volumes) > 0 {
				var volumesList []string
				for volume := range volumes {
					volumesList = append(volumesList, volume)
				}
				addMultiLineData(&rows, "Volumes", volumesList, "", 0)
			}

			// Entrypoint
			if entrypoint, ok := config["Entrypoint"].([]any); ok && len(entrypoint) > 0 {
				var entryList []string
				for _, e := range entrypoint {
					entryList = append(entryList, fmt.Sprintf("%v", e))
				}
				addRow(&rows, "Entrypoint", strings.Join(entryList, " "))
			}

			// Command
			if cmd, ok := config["Cmd"].([]any); ok && len(cmd) > 0 {
				var cmdList []string
				for _, c := range cmd {
					cmdList = append(cmdList, fmt.Sprintf("%v", c))
				}
				addRow(&rows, "Command", strings.Join(cmdList, " "))
			}

			// Labels
			if labels, ok := config["Labels"].(map[string]any); ok && len(labels) > 0 {
				var labelsList []string
				for k, v := range labels {
					labelsList = append(labelsList, fmt.Sprintf("%s: %v", k, v))
				}
				addMultiLineData(&rows, "Labels", labelsList, "", 0)
			}
		}
	}

	// SECTION: Other Attributes
	if artifact.ExtraAttrs != nil {
		otherAttrs := false

		for key, value := range artifact.ExtraAttrs {
			if key == "architecture" || key == "os" || key == "config" || key == "created" || key == "layers" {
				continue
			}

			if !otherAttrs {
				addSectionHeader(&rows, "OTHER ATTRIBUTES")
				otherAttrs = true
			}

			switch v := value.(type) {
			case string, bool, int, int64, float64:
				addRow(&rows, key, fmt.Sprintf("%v", v))
			default:
				// For complex structures, pretty print as JSON
				jsonBytes, err := json.MarshalIndent(value, "", "  ")
				if err == nil {
					jsonString := string(jsonBytes)
					lines := strings.Split(jsonString, "\n")
					addMultiLineData(&rows, key, lines, "", 0)
				} else {
					addRow(&rows, key, "(complex structure - see JSON output)")
				}
			}
		}
	}

	// SECTION: Vulnerabilities
	addSectionHeader(&rows, "SECURITY")
	var totalVulnerabilities int64
	for _, scan := range artifact.ScanOverview {
		totalVulnerabilities += scan.Summary.Total
	}
	addRow(&rows, "Vulnerabilities", strconv.FormatInt(totalVulnerabilities, 10))

	// SECTION: References
	if len(artifact.References) > 0 {
		addSectionHeader(&rows, "REFERENCES")
		for i, ref := range artifact.References {
			if ref.Platform != nil && ref.Platform.Architecture != "" {
				addRow(&rows, fmt.Sprintf("Reference %d", i+1), fmt.Sprintf("%s: %s", ref.Platform.Architecture, ref.ChildDigest))
			}
		}
	}

	// SECTION: Additional links
	if len(artifact.AdditionLinks) > 0 {
		addSectionHeader(&rows, "ADDITIONAL INFO")
		var links []string
		for key := range artifact.AdditionLinks {
			links = append(links, key)
		}
		addRow(&rows, "Available Endpoints", strings.Join(links, ", "))
	}

	// Display the table
	m := tablelist.NewModel(detailsColumns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error displaying artifact details:", err)
		os.Exit(1)
	}
}

func addSectionHeader(rows *[]table.Row, title string) {
	*rows = append(*rows, table.Row{"", ""})
	*rows = append(*rows, table.Row{"--- " + title + " ---", ""})
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
