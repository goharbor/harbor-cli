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
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/harbor-cli/pkg/utils"
	helpers "github.com/goharbor/harbor-cli/test/helper"
	"github.com/stretchr/testify/assert"
)

func setupTestEnvironment(t *testing.T) (*utils.HarborData, error) {
	tempDir := t.TempDir()
	helpers.SetMockKeyring(t)
	data := helpers.Initialize(t, tempDir)

	err := utils.GenerateEncryptionKey()
	assert.NoError(t, err)
	key, err := utils.GetEncryptionKey()
	assert.NoError(t, err)
	encrypted, err := utils.Encrypt(key, []byte("test-password"))
	assert.NoError(t, err)
	testConfig := &utils.HarborConfig{
		Credentials: []utils.Credential{
			{
				Name:          "test-credential",
				Username:      "test-user",
				Password:      encrypted,
				ServerAddress: "https://test-endpoint",
			},
		},
		CurrentCredentialName: "test-credential",
	}
	err = utils.SetCurrentHarborConfig(testConfig)
	assert.NoError(t, err)
	return data, err
}

func TestGetClient_Success(t *testing.T) {
	data, err := setupTestEnvironment(t)
	assert.NoError(t, err)
	defer helpers.ConfigCleanup(t, data)
	client, getErr := utils.GetClient()
	assert.NoError(t, getErr)
	assert.NotNil(t, client)
}

func TestContextWithClient_Success(t *testing.T) {
	data, err := setupTestEnvironment(t)
	assert.NoError(t, err)
	defer helpers.ConfigCleanup(t, data)
	ctx, client, err := utils.ContextWithClient()
	assert.NotNil(t, ctx)
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestContextWithClient_Failure(t *testing.T) {
	ctx, client, err := utils.ContextWithClient()
	assert.Nil(t, ctx)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestGetClientByConfig(t *testing.T) {
	clientConfig := &harbor.ClientSetConfig{
		URL:      "https://test-endpoint",
		Username: "test-user",
		Password: "test-password",
	}
	client := utils.GetClientByConfig(clientConfig)
	assert.NotNil(t, client)
}

func TestGetClientByCredentialName(t *testing.T) {
	data, err := setupTestEnvironment(t)
	assert.NoError(t, err)
	defer helpers.ConfigCleanup(t, data)
	client, clientErr := utils.GetClientByCredentialName("test-credential")
	assert.NotNil(t, client)
	assert.NoError(t, clientErr)
}

func TestGetClientByCredentialName_Failure(t *testing.T) {
	data, err := setupTestEnvironment(t)
	assert.NoError(t, err)
	defer helpers.ConfigCleanup(t, data)
	client, clientErr := utils.GetClientByCredentialName("non-existent-credential")
	assert.Nil(t, client)
	assert.Error(t, clientErr)
}

func TestGetClient_EmptyCredentialName(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SetMockKeyring(t)
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	// Set config with empty credential name
	testConfig := &utils.HarborConfig{
		Credentials:           []utils.Credential{},
		CurrentCredentialName: "",
	}
	err := utils.SetCurrentHarborConfig(testConfig)
	assert.NoError(t, err)

	// Reset ClientOnce to allow re-initialization
	utils.ClientOnce = sync.Once{}
	utils.ClientInstance = nil
	utils.ClientErr = nil

	client, getErr := utils.GetClient()
	assert.Nil(t, client)
	assert.Error(t, getErr)
	assert.Contains(t, getErr.Error(), "current-credential-name is not set in config file")
}

func TestGetClient_InvalidCredentialName(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SetMockKeyring(t)
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	err := utils.GenerateEncryptionKey()
	assert.NoError(t, err)

	// Set config with invalid credential name
	testConfig := &utils.HarborConfig{
		Credentials: []utils.Credential{
			{
				Name:          "valid-credential",
				Username:      "test-user",
				Password:      "encrypted-pass",
				ServerAddress: "https://test-endpoint",
			},
		},
		CurrentCredentialName: "invalid-credential",
	}
	err = utils.SetCurrentHarborConfig(testConfig)
	assert.NoError(t, err)

	// Reset ClientOnce to allow re-initialization
	utils.ClientOnce = sync.Once{}
	utils.ClientInstance = nil
	utils.ClientErr = nil

	client, getErr := utils.GetClient()
	assert.Nil(t, client)
	assert.Error(t, getErr)
	assert.Contains(t, getErr.Error(), "failed to get credential")
}

func TestGetClientByCredentialName_DecryptionError(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SetMockKeyring(t)
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	err := utils.GenerateEncryptionKey()
	assert.NoError(t, err)

	// Create credential with invalid encrypted password
	testConfig := &utils.HarborConfig{
		Credentials: []utils.Credential{
			{
				Name:          "test-credential",
				Username:      "test-user",
				Password:      "invalid-encrypted-data",
				ServerAddress: "https://test-endpoint",
			},
		},
		CurrentCredentialName: "test-credential",
	}
	err = utils.SetCurrentHarborConfig(testConfig)
	assert.NoError(t, err)

	client, clientErr := utils.GetClientByCredentialName("test-credential")
	assert.Nil(t, client)
	assert.Error(t, clientErr)
	assert.Contains(t, clientErr.Error(), "failed to decrypt password")
}