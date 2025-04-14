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

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: tablelist.WidthXS},
	{Title: "Artifact Digest", Width: tablelist.WidthL},
	{Title: "Type", Width: tablelist.WidthS},
	{Title: "Size", Width: tablelist.WidthS},
	{Title: "Vulnerabilities", Width: tablelist.WidthM},
	{Title: "Push Time", Width: tablelist.WidthM},
}

func ViewArtifact(artifact *models.Artifact) {
	var rows []table.Row

	pushTime, _ := utils.FormatCreatedTime(artifact.PushTime.String())
	artifactSize := utils.FormatSize(artifact.Size)
	var totalVulnerabilities int64
	for _, scan := range artifact.ScanOverview {
		totalVulnerabilities += scan.Summary.Total
	}
	rows = append(rows, table.Row{
		strconv.FormatInt(int64(artifact.ID), 10),
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
