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
package list

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestCreateSystemInfo_WithVolumes(t *testing.T) {
	generalInfo := &models.GeneralInfo{}
	stats := &models.Statistic{}
	volumes := &models.SystemInfo{
		Storage: []*models.Storage{
			{
				Free:  1000,
				Total: 5000,
			},
		},
	}
	cliInfo := &api.CLIInfo{
		Username:           "admin",
		RegistryAddress:    "localhost",
		IsSysAdmin:         true,
		PreviouslyLoggedIn: []string{},
	}

	result := CreateSystemInfo(generalInfo, stats, volumes, cliInfo, "v1.0.0", "linux")

	assert.NotNil(t, result.VolumeInfo)
	assert.Equal(t, uint64(1000), result.VolumeInfo.Free)
	assert.Equal(t, uint64(5000), result.VolumeInfo.Total)
}

func TestCreateSystemInfo_NilVolumes(t *testing.T) {
	generalInfo := &models.GeneralInfo{}
	stats := &models.Statistic{}
	var volumes *models.SystemInfo // Simulate missing payload

	cliInfo := &api.CLIInfo{
		Username:           "admin",
		RegistryAddress:    "localhost",
		IsSysAdmin:         true,
		PreviouslyLoggedIn: []string{},
	}

	result := CreateSystemInfo(generalInfo, stats, volumes, cliInfo, "v1.0.0", "linux")

	assert.NotNil(t, result.VolumeInfo)
	assert.Equal(t, uint64(0), result.VolumeInfo.Free)
	assert.Equal(t, uint64(0), result.VolumeInfo.Total)
}
