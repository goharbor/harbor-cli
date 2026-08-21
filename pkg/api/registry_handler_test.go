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
	"errors"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestGetRegistryResponse_ViewRegistryError(t *testing.T) {
	// Save original and restore after test
	originalViewRegistry := viewRegistryFunc
	defer func() { viewRegistryFunc = originalViewRegistry }()

	// Mock ViewRegistry to return an error
	viewRegistryFunc = func(registryId int64) (*registry.GetRegistryOK, error) {
		return nil, errors.New("connection refused")
	}

	reg, err := GetRegistryResponse(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
	assert.Nil(t, reg)
}

func TestGetRegistryResponse_Success(t *testing.T) {
	originalViewRegistry := viewRegistryFunc
	defer func() { viewRegistryFunc = originalViewRegistry }()

	// Mock ViewRegistry to return a registry
	viewRegistryFunc = func(registryId int64) (*registry.GetRegistryOK, error) {
		return &registry.GetRegistryOK{
			Payload: &models.Registry{ID: 1, Name: "dockerhub"},
		}, nil
	}

	reg, err := GetRegistryResponse(1)

	assert.NoError(t, err)
	assert.NotNil(t, reg)
	assert.Equal(t, int64(1), reg.ID)
	assert.Equal(t, "dockerhub", reg.Name)
}

func TestGetRegistryIdByName_ListRegistriesError(t *testing.T) {
	originalListRegistries := listRegistriesFunc
	defer func() { listRegistriesFunc = originalListRegistries }()

	// Mock ListRegistries to return an error
	listRegistriesFunc = func(opts ...ListFlags) (*registry.ListRegistriesOK, error) {
		return nil, errors.New("connection refused")
	}

	id, err := GetRegistryIdByName("dockerhub")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
	assert.Equal(t, int64(0), id)
}

func TestGetRegistryIdByName_NotFound(t *testing.T) {
	originalListRegistries := listRegistriesFunc
	defer func() { listRegistriesFunc = originalListRegistries }()

	// Mock ListRegistries to return a list without the requested registry
	listRegistriesFunc = func(opts ...ListFlags) (*registry.ListRegistriesOK, error) {
		return &registry.ListRegistriesOK{
			Payload: []*models.Registry{
				{ID: 1, Name: "dockerhub"},
				{ID: 2, Name: "quay"},
			},
		}, nil
	}

	id, err := GetRegistryIdByName("missing-registry")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'missing-registry' not found")
	assert.Equal(t, int64(0), id)
}

func TestGetRegistryIdByName_Success(t *testing.T) {
	originalListRegistries := listRegistriesFunc
	defer func() { listRegistriesFunc = originalListRegistries }()

	// Mock ListRegistries to return a list containing the requested registry
	listRegistriesFunc = func(opts ...ListFlags) (*registry.ListRegistriesOK, error) {
		return &registry.ListRegistriesOK{
			Payload: []*models.Registry{
				{ID: 1, Name: "dockerhub"},
				{ID: 2, Name: "quay"},
			},
		}, nil
	}

	id, err := GetRegistryIdByName("quay")

	assert.NoError(t, err)
	assert.Equal(t, int64(2), id)
}
