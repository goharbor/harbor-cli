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
package utils_test

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractConfigValues_FromConfigurations(t *testing.T) {
	authMode := "db_auth"
	selfReg := true
	tokenExp := int64(30)

	cfg := &models.Configurations{
		AuthMode:         &authMode,
		SelfRegistration: &selfReg,
		TokenExpiration:  &tokenExp,
	}

	result := utils.ExtractConfigValues(cfg)

	assert.Equal(t, "db_auth", result["AuthMode"])
	assert.Equal(t, true, result["SelfRegistration"])
	assert.Equal(t, int64(30), result["TokenExpiration"])
}

func TestExtractConfigValues_FromConfigurationsResponse(t *testing.T) {
	resp := &models.ConfigurationsResponse{
		AuthMode:         &models.StringConfigItem{Value: "ldap_auth"},
		SelfRegistration: &models.BoolConfigItem{Value: true},
		TokenExpiration:  &models.IntegerConfigItem{Value: int64(60)},
	}

	result := utils.ExtractConfigValues(resp)

	assert.Equal(t, "ldap_auth", result["AuthMode"])
	assert.Equal(t, true, result["SelfRegistration"])
	assert.Equal(t, int64(60), result["TokenExpiration"])
}

func TestExtractConfigValues_Nil(t *testing.T) {
	result := utils.ExtractConfigValues[*models.Configurations](nil)
	assert.Empty(t, result)
}

func TestExtractConfigValues_ExcludesEmptyString(t *testing.T) {
	authMode := ""
	cfg := &models.Configurations{AuthMode: &authMode}

	result := utils.ExtractConfigValues(cfg)
	_, exists := result["AuthMode"]
	assert.False(t, exists)
}

func TestExtractConfigValues_StringPointerCase(t *testing.T) {
	mode := "ldap_auth"
	cfg := &models.Configurations{AuthMode: &mode}

	result := utils.ExtractConfigValues(cfg)
	assert.Equal(t, "ldap_auth", result["AuthMode"])
}

func TestExtractConfigValues_GenericPointerFallback(t *testing.T) {
	quota := int64(100)
	enabled := true
	cfg := &models.Configurations{StoragePerProject: &quota, SelfRegistration: &enabled}

	result := utils.ExtractConfigValues(cfg)
	assert.Equal(t, int64(100), result["StoragePerProject"])
	assert.Equal(t, true, result["SelfRegistration"])
}

func TestConvertToConfigurations(t *testing.T) {
	resp := &models.ConfigurationsResponse{
		AuthMode:         &models.StringConfigItem{Value: "db_auth"},
		SelfRegistration: &models.BoolConfigItem{Value: true},
		TokenExpiration:  &models.IntegerConfigItem{Value: int64(30)},
	}

	cfg := utils.ConvertToConfigurations(resp)

	assert.NotNil(t, cfg)
	assert.Equal(t, "db_auth", *cfg.AuthMode)
	assert.Equal(t, true, *cfg.SelfRegistration)
	assert.Equal(t, int64(30), *cfg.TokenExpiration)
}

func TestConvertToConfigurations_Nil(t *testing.T) {
	cfg := utils.ConvertToConfigurations(nil)
	assert.NotNil(t, cfg)
	assert.Equal(t, &models.Configurations{}, cfg)
}

func TestConvertToConfigurations_SecretFieldEncrypted(t *testing.T) {
	utils.SetKeyringProvider(&utils.EnvironmentKeyring{EnvVarName: "TEST_REFLECT_ENC_KEY"})
	t.Cleanup(func() { utils.SetKeyringProvider(&utils.SystemKeyring{}) })
	require.NoError(t, utils.GenerateEncryptionKey())

	resp := &models.ConfigurationsResponse{
		UaaClientSecret: &models.StringConfigItem{Value: "super-secret"},
		AuthMode:        &models.StringConfigItem{Value: "db_auth"},
	}
	cfg := utils.ConvertToConfigurations(resp)

	require.NotNil(t, cfg.UaaClientSecret)
	assert.NotEqual(t, "super-secret", *cfg.UaaClientSecret)
	key, err := utils.GetEncryptionKey()
	require.NoError(t, err)
	plain, err := utils.Decrypt(key, *cfg.UaaClientSecret)
	require.NoError(t, err)
	assert.Equal(t, "super-secret", plain)

	require.NotNil(t, cfg.AuthMode)
	assert.Equal(t, "db_auth", *cfg.AuthMode)
}

func TestConvertToConfigurations_SecretFieldEncryptError(t *testing.T) {
	utils.SetKeyringProvider(&utils.EnvironmentKeyring{EnvVarName: "TEST_REFLECT_ENC_BAD_KEY"})
	t.Cleanup(func() { utils.SetKeyringProvider(&utils.SystemKeyring{}) })
	t.Setenv("TEST_REFLECT_ENC_BAD_KEY", "!!!not-base64!!!")

	resp := &models.ConfigurationsResponse{
		UaaClientSecret: &models.StringConfigItem{Value: "super-secret"},
		AuthMode:        &models.StringConfigItem{Value: "db_auth"},
	}
	cfg := utils.ConvertToConfigurations(resp)

	assert.Nil(t, cfg.UaaClientSecret)
	require.NotNil(t, cfg.AuthMode)
	assert.Equal(t, "db_auth", *cfg.AuthMode)
}

func TestExtractConfigurationsByCategory_Nil(t *testing.T) {
	cfg := utils.ExtractConfigurationsByCategory(nil, "authentication")
	assert.NotNil(t, cfg)
}

func TestExtractConfigurationsByCategory_FiltersByCategory(t *testing.T) {
	resp := &models.ConfigurationsResponse{
		AuthMode:                   &models.StringConfigItem{Value: "db_auth"},
		SelfRegistration:           &models.BoolConfigItem{Value: true},
		ProjectCreationRestriction: &models.StringConfigItem{Value: "adminonly"},
	}

	cfg := utils.ExtractConfigurationsByCategory(resp, "authentication")

	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.AuthMode)
	assert.Equal(t, "db_auth", *cfg.AuthMode)
	assert.Nil(t, cfg.SelfRegistration)
	assert.Nil(t, cfg.ProjectCreationRestriction)
}

func TestIsCategory(t *testing.T) {
	assert.True(t, utils.IsCategory("AuthMode", "authentication"))
	assert.False(t, utils.IsCategory("AuthMode", "security"))
	assert.False(t, utils.IsCategory("Nonexistent", "authentication"))
	assert.True(t, utils.IsCategory("Anything", ""))
}
