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
package scanner

import (
	"errors"
	"testing"

	client_scanner "github.com/goharbor/go-client/pkg/sdk/v2.0/client/scanner"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCommand_NonInteractiveWithFlags(t *testing.T) {
	oldGetScannerByName := getScannerByNameFunc
	oldGetScanner := getScannerFunc
	oldUpdateScannerView := updateScannerViewFunc
	oldUpdateScanner := updateScannerFunc
	defer func() {
		getScannerByNameFunc = oldGetScannerByName
		getScannerFunc = oldGetScanner
		updateScannerViewFunc = oldUpdateScannerView
		updateScannerFunc = oldUpdateScanner
	}()

	mockScanner := models.ScannerRegistration{
		UUID:        "test-uuid-123",
		Name:        "old-scanner",
		Description: "old description",
		URL:         "http://old.example.com",
	}

	getScannerByNameFunc = func(name string) (models.ScannerRegistration, error) {
		assert.Equal(t, "old-scanner", name)
		return mockScanner, nil
	}

	getScannerFunc = func(registrationID string) (client_scanner.GetScannerOK, error) {
		assert.Equal(t, "test-uuid-123", registrationID)
		return client_scanner.GetScannerOK{
			Payload: &mockScanner,
		}, nil
	}

	viewCalled := false
	updateScannerViewFunc = func(s *models.ScannerRegistration) error {
		viewCalled = true
		return nil
	}

	var updatedRegistration models.ScannerRegistration
	updateScannerFunc = func(registrationID string, s models.ScannerRegistration) error {
		assert.Equal(t, "test-uuid-123", registrationID)
		updatedRegistration = s
		return nil
	}

	cmd := UpdateCommand()
	cmd.SetArgs([]string{"old-scanner", "--name", "new-scanner", "--description", "new description"})
	err := cmd.Execute()

	assert.NoError(t, err)
	assert.False(t, viewCalled, "interactive view should NOT be called when flags are provided")
	assert.Equal(t, "new-scanner", updatedRegistration.Name)
	assert.Equal(t, "new description", updatedRegistration.Description)
}

func TestUpdateCommand_InteractiveWhenNoFlags(t *testing.T) {
	oldGetScannerByName := getScannerByNameFunc
	oldGetScanner := getScannerFunc
	oldUpdateScannerView := updateScannerViewFunc
	oldUpdateScanner := updateScannerFunc
	defer func() {
		getScannerByNameFunc = oldGetScannerByName
		getScannerFunc = oldGetScanner
		updateScannerViewFunc = oldUpdateScannerView
		updateScannerFunc = oldUpdateScanner
	}()

	mockScanner := models.ScannerRegistration{
		UUID:        "test-uuid-456",
		Name:        "scanner-two",
		Description: "description two",
		URL:         "http://two.example.com",
	}

	getScannerByNameFunc = func(name string) (models.ScannerRegistration, error) {
		return mockScanner, nil
	}

	getScannerFunc = func(registrationID string) (client_scanner.GetScannerOK, error) {
		return client_scanner.GetScannerOK{
			Payload: &mockScanner,
		}, nil
	}

	viewCalled := false
	updateScannerViewFunc = func(s *models.ScannerRegistration) error {
		viewCalled = true
		s.Name = "updated-via-tui"
		return nil
	}

	var updatedRegistration models.ScannerRegistration
	updateScannerFunc = func(registrationID string, s models.ScannerRegistration) error {
		updatedRegistration = s
		return nil
	}

	cmd := UpdateCommand()
	cmd.SetArgs([]string{"scanner-two"})
	err := cmd.Execute()

	assert.NoError(t, err)
	assert.True(t, viewCalled, "interactive view SHOULD be called when no flags are provided")
	assert.Equal(t, "updated-via-tui", updatedRegistration.Name)
}

func TestUpdateCommand_NoScannerSelected(t *testing.T) {
	oldGetScannerIdFromUser := getScannerIdFromUserFunc
	defer func() {
		getScannerIdFromUserFunc = oldGetScannerIdFromUser
	}()

	getScannerIdFromUserFunc = func() string {
		return ""
	}

	cmd := UpdateCommand()
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scanner ID cannot be empty")
}

func TestUpdateCommand_ScannerNotFound(t *testing.T) {
	oldGetScannerByName := getScannerByNameFunc
	defer func() {
		getScannerByNameFunc = oldGetScannerByName
	}()

	getScannerByNameFunc = func(name string) (models.ScannerRegistration, error) {
		return models.ScannerRegistration{}, errors.New("not found")
	}

	cmd := UpdateCommand()
	cmd.SetArgs([]string{"non-existent"})
	err := cmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to retrieve scanner by name")
}
