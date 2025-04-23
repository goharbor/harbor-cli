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
	"fmt"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/preheat"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/instance/create"
	log "github.com/sirupsen/logrus"
)

func CreateInstance(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Preheat.CreateInstance(ctx, &preheat.CreateInstanceParams{Instance: &models.Instance{Vendor: opts.Vendor, Name: opts.Name, Description: opts.Description, Endpoint: strings.TrimSpace(opts.Endpoint), Enabled: opts.Enabled, AuthMode: opts.AuthMode, AuthInfo: opts.AuthInfo, Insecure: opts.Insecure}})
	if err != nil {
		return err
	}

	log.Infof("Instance %s created", opts.Name)
	return nil
}

func DeleteInstance(instanceName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Preheat.DeleteInstance(ctx, &preheat.DeleteInstanceParams{PreheatInstanceName: instanceName})

	if err != nil {
		return err
	}

	log.Info("instance deleted successfully")

	return nil
}

func ListInstance(opts ...ListFlags) (*preheat.ListInstancesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Preheat.ListInstances(ctx, &preheat.ListInstancesParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Sort:     &listFlags.Sort,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetInstanceNameByID(id int64) (string, error) {
	instances, err := ListInstance()
	if err != nil {
		return "", err
	}

	for _, inst := range instances.Payload {
		if inst.ID == id {
			return inst.Name, nil
		}
	}

	return "", fmt.Errorf("no instance found with ID: %v", id)
}
