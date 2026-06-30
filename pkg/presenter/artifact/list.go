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
package presenterartifact

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
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

func ListArtifacts(projectName, repoName string, opts api.ListFlags) error {
	m := tablelistv2.NewModel(columns, LoadArtifactList(projectName, repoName, opts))

	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	if modelErr := finalModel.(tablelistv2.Model).Error; modelErr != nil {
		return modelErr
	}

	return nil
}

func LoadArtifactList(project, repo string, listOpts api.ListFlags) func() ([]table.Row, error) {
	projectName := project
	repoName := repo
	opts := listOpts

	return func() ([]table.Row, error) {
		artifacts, err := api.ListArtifact(projectName, repoName, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list artifacts: %v", err)
		}

		var rows []table.Row
		for _, artifact := range artifacts.Payload {
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
		}

		return rows, nil
	}
}
