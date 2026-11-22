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
package artifact

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
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
	model := tablelistv2.NewModel(
		columns,
		fmt.Sprintf("Loading artifacts for %s/%s...", projectName, repoName),
		loadArtifactRows(projectName, repoName, opts),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running artifact list: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected artifact list model result")
	}

	return loadedModel.Err
}

func loadArtifactRows(projectName, repoName string, opts api.ListFlags) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		artifacts, err := api.ListArtifact(projectName, repoName, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list artifacts: %v", err)
		}

		return buildArtifactRows(artifacts.Payload), nil
	}
}

func buildArtifactRows(artifacts []*models.Artifact) []table.Row {
	rows := make([]table.Row, 0, len(artifacts))

	for _, artifact := range artifacts {
		rows = append(rows, table.Row{
			strconv.FormatInt(int64(artifact.ID), 10),
			tagNames(artifact),
			shortDigest(artifact.Digest),
			artifact.Type,
			formatArtifactSize(artifact),
			totalVulnerabilities(artifact),
			formatPushTime(artifact),
		})
	}

	return rows
}

func tagNames(artifact *models.Artifact) string {
	var names []string
	for _, tag := range artifact.Tags {
		names = append(names, tag.Name)
	}
	if len(names) == 0 {
		return "-"
	}
	return strings.Join(names, ", ")
}

func shortDigest(digest string) string {
	if len(digest) <= 16 {
		return digest
	}
	return digest[:16]
}

func formatArtifactSize(artifact *models.Artifact) string {
	return utils.FormatSize(artifact.Size)
}

func totalVulnerabilities(artifact *models.Artifact) string {
	var total int64
	for _, scan := range artifact.ScanOverview {
		total += scan.Summary.Total
	}
	return strconv.FormatInt(total, 10)
}

func formatPushTime(artifact *models.Artifact) string {
	formattedTime, err := utils.FormatCreatedTime(artifact.PushTime.String())
	if err != nil {
		return artifact.PushTime.String()
	}
	return formattedTime
}
