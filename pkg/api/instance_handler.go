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
	"strconv"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/preheat"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/instance/create"
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

	fmt.Printf("Instance %s created\n", opts.Name)
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

	fmt.Println("Instance deleted successfully")

	return nil
}

func UpdateInstance(instanceName string, instance models.Instance) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	instance.Endpoint = strings.TrimSpace(instance.Endpoint)

	_, err = client.Preheat.UpdateInstance(ctx, &preheat.UpdateInstanceParams{
		PreheatInstanceName: instanceName,
		Instance:            &instance,
	})
	if err != nil {
		return err
	}

	log.Infof("Instance %s updated", instance.Name)
	return nil
}

func PingInstance(instanceNameOrID string, useInstanceID bool) (*preheat.PingInstancesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	instance, err := GetInstance(instanceNameOrID, useInstanceID)
	if err != nil {
		return nil, err
	}
	if instance == nil || instance.Payload == nil {
		return nil, fmt.Errorf("failed to ping instance: empty response")
	}

	response, err := client.Preheat.PingInstances(ctx, &preheat.PingInstancesParams{
		Instance: instance.Payload,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ListAllInstance(opts ...ListFlags) (*preheat.ListInstancesOK, error) {
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
	instances, err := ListAllInstance()
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

func GetInstance(instanceNameOrID string, useInstanceID bool) (*preheat.GetInstanceOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	if useInstanceID {
		instanceID, err := strconv.ParseInt(instanceNameOrID, 10, 64)
		if err != nil {
			return nil, err
		}
		instanceNameOrID, err = GetInstanceNameByID(instanceID)
		if err != nil {
			return nil, err
		}
	}

	response, err := client.Preheat.GetInstance(ctx, &preheat.GetInstanceParams{
		PreheatInstanceName: instanceNameOrID,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
