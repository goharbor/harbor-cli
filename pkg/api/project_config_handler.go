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
package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project_metadata"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"

	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func ListConfig(isID bool, projectNameOrID string) (*project_metadata.ListProjectMetadatasOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	isName := !isID
	response, err := client.ProjectMetadata.ListProjectMetadatas(ctx, &project_metadata.ListProjectMetadatasParams{ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func UpdateConfig(isID bool, projectNameOrID string, metadata models.ProjectMetadata) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	isName := !isID
	response, err := client.Project.UpdateProject(ctx, &project.UpdateProjectParams{
		ProjectNameOrID: projectNameOrID,
		XIsResourceName: &isName,
		Project: &models.ProjectReq{
			Metadata: &metadata,
		},
	})
	if err != nil {
		return err
	}
	if response != nil {
		log.Info("Metadata updated successfully")
	}

	return nil
}
