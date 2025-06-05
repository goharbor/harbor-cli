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

var columns = []table.Column{
	{Title: "ID", Width: tablelist.WidthS},
	{Title: "Tags", Width: tablelist.WidthL},
	{Title: "Artifact Digest", Width: tablelist.WidthXL},
	{Title: "Type", Width: tablelist.WidthS},
	{Title: "Size", Width: tablelist.WidthM},
	{Title: "Vulnerabilities", Width: tablelist.WidthL},
	{Title: "Push Time", Width: tablelist.WidthL},
}

func ViewArtifact(artifact *models.Artifact) {
	showBasicInfo(artifact)
	//for extra details
	fmt.Println("\n[Artifact Details]")
	showExtraInfo(artifact)
	//show addon data if available
	if artifact.AdditionLinks != nil {
		fmt.Println("\n[Additional Information]")
		for key := range artifact.AdditionLinks {
			fmt.Printf("- %s available\n", key)
		}
	}
}

func showBasicInfo(artifact *models.Artifact) {
	var rows []table.Row

	pushTime, _ := utils.FormatCreatedTime(artifact.PushTime.String())
	artifactSize := utils.FormatSize(artifact.Size)
	var tagNames []string
	for _, tag := range artifact.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	tags := "-"
	if len(tagNames) > 0 {
		tags = strings.Join(tagNames, ", ")
	}

	var totalVulnerabilities int64
	for _, scan := range artifact.ScanOverview {
		totalVulnerabilities += scan.Summary.Total
	}
	rows = append(rows, table.Row{
		strconv.FormatInt(int64(artifact.ID), 10),
		tags,
		artifact.Digest[:16],
		artifact.Type,
		artifactSize,
		strconv.FormatInt(totalVulnerabilities, 10),
		pushTime,
	})

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func showExtraInfo(artifact *models.Artifact) {
	fmt.Printf("RepositoryID: %d\n", artifact.RepositoryID)
	fmt.Printf("Media Type: %s\n", artifact.MediaType)

	if artifact.ExtraAttrs != nil {
		var configSection map[string]any
		if arch, ok := artifact.ExtraAttrs["architecture"].(string); ok {
			fmt.Printf("Architecture: %s\n", arch)
		}

		if os, ok := artifact.ExtraAttrs["os"].(string); ok {
			fmt.Printf("OS: %s\n", os)
		}

		if config, ok := artifact.ExtraAttrs["config"].(map[string]any); ok {
			configSection = config
			if author, ok := config["author"].(string); ok {
				fmt.Printf("Author: %s\n", author)
			}
		}

		if created, ok := artifact.ExtraAttrs["created"].(string); ok {
			fmt.Printf("Created: %s\n", created)
		}

		if layers, ok := artifact.ExtraAttrs["layers"].([]any); ok {
			fmt.Printf("Layers: %d\n", len(layers))
		}

		if configSection != nil {
			fmt.Println("\n[Config Details]")

			//for env variables if available
			if env, ok := configSection["Env"].([]any); ok && len(env) > 0 {
				fmt.Println("Environment Variables:")
				for _, e := range env {
					fmt.Printf("  - %s\n", e)
				}
			}

			//for exposed ports if available
			if ports, ok := configSection["ExposedPorts"].(map[string]any); ok {
				fmt.Println("Exposed Ports:")
				for port := range ports {
					fmt.Printf("  - %s\n", port)
				}
			}

			//for volumes if available
			if volumes, ok := configSection["Volumes"].(map[string]any); ok {
				fmt.Println("Volumes:")
				for volume := range volumes {
					fmt.Printf("  - %s\n", volume)
				}
			}

			if entrypoint, ok := configSection["Entrypoint"].([]any); ok {
				fmt.Printf("Entrypoint: %v\n", entrypoint)
			}

			if cmd, ok := configSection["Cmd"].([]any); ok {
				fmt.Printf("Command: %v\n", cmd)
			}
		}

		//for labels in config
		if configSection != nil {
			if labels, ok := configSection["Labels"].(map[string]any); ok && len(labels) > 0 {
				fmt.Println("\n[Labels]")
				for key, value := range labels {
					fmt.Printf("%s: %v\n", key, value)
				}
			}
		}

		//for other interesting fields
		fmt.Println("\n[Other Attributes]")
		for key, value := range artifact.ExtraAttrs {
			//skipping already displayed ones
			if key == "architecture" || key == "os" || key == "config" ||
				key == "created" || key == "layers" {
				continue
			}

			switch v := value.(type) {
			case string, bool, int, int64, float64:
				fmt.Printf("%s: %v\n", key, v)
			default:
				fmt.Printf("%s: (complex data available)\n", key)
			}
		}
	}

	if len(artifact.References) > 0 {
		fmt.Println("\n[References]")
		for _, ref := range artifact.References {
			fmt.Printf("%s: %s\n", ref.Platform.Architecture, ref.ChildDigest)
		}
	}
}
