package utils_test

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
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

func TestExtractConfigurationsByCategory_Nil(t *testing.T) {
	cfg := utils.ExtractConfigurationsByCategory(nil, "authentication")
	assert.NotNil(t, cfg)
}

func TestExtractConfigurationsByCategory_FiltersByCategory(t *testing.T) {
	resp := &models.ConfigurationsResponse{
		AuthMode:                  &models.StringConfigItem{Value: "db_auth"},
		SelfRegistration:          &models.BoolConfigItem{Value: true},
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
