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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/statistic"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetStats() (*statistic.GetStatisticOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Statistic.GetStatistic(
		ctx,
		&statistic.GetStatisticParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetSystemInfo() (*systeminfo.GetSystemInfoOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Systeminfo.GetSystemInfo(
		ctx,
		&systeminfo.GetSystemInfoParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetSystemVolumes() (*systeminfo.GetVolumesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Systeminfo.GetVolumes(
		ctx,
		&systeminfo.GetVolumesParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
