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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/configure"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetConfigurations() (*configure.GetConfigurationsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Configure.GetConfigurations(ctx, &configure.GetConfigurationsParams{})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func UpdateConfigurations(config *utils.HarborConfig) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	params := &configure.UpdateConfigurationsParams{
		Configurations: &config.Configurations,
	}

	_, err = client.Configure.UpdateConfigurations(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
