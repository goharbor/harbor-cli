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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/statistic"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/viper"
)

type CLIInfo struct {
	Username           string
	RegistryAddress    string
	IsSysAdmin         bool
	PreviouslyLoggedIn []string
	OSinfo             string
}

func GetStats() (*statistic.GetStatisticOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("get stats: %w", err)
	}
	return client.Statistic.GetStatistic(ctx, &statistic.GetStatisticParams{})
}

func GetSystemInfo() (*systeminfo.GetSystemInfoOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("get system info: %w", err)
	}
	return client.Systeminfo.GetSystemInfo(ctx, &systeminfo.GetSystemInfoParams{})
}

func GetSystemVolumes() (*systeminfo.GetVolumesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("get system volumes: %w", err)
	}
	return client.Systeminfo.GetVolumes(ctx, &systeminfo.GetVolumesParams{})
}

func GetCLIInfo() (*CLIInfo, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("cli info: failed to create Harbor client: %w", err)
	}

	currentCred := viper.GetString("current-credential-name")
	if currentCred == "" {
		return nil, fmt.Errorf("cli info: no active credentials found")
	}

	creds, ok := viper.Get("credentials").([]interface{})
	if !ok {
		return nil, fmt.Errorf("cli info: invalid type for credentials, expected []interface{}")
	}
	var registryAddress string
	seen := make(map[string]struct{})
	var previousRegistries []string

	for _, cred := range creds {
		c := cred.(map[string]interface{})
		addr := c["serveraddress"].(string)
		name := c["name"].(string)

		if name == currentCred {
			registryAddress = addr
		}
		if _, exists := seen[addr]; !exists {
			previousRegistries = append(previousRegistries, name)
			seen[addr] = struct{}{}
		}
	}

	userResp, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
	if err != nil {
		return nil, fmt.Errorf("cli info: get current user info: %w", err)
	}

	return &CLIInfo{
		Username:           userResp.Payload.Username,
		RegistryAddress:    registryAddress,
		IsSysAdmin:         userResp.Payload.SysadminFlag,
		PreviouslyLoggedIn: previousRegistries,
	}, nil
}
