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
package update

import (
	"errors"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	assert.Error(t, ValidateEmail(""))
	assert.Error(t, ValidateEmail("   "))
	assert.Error(t, ValidateEmail("invalid-email"))
	assert.NoError(t, ValidateEmail("test@example.com"))
}

func TestValidateRealname(t *testing.T) {
	assert.Error(t, ValidateRealname(""))
	assert.Error(t, ValidateRealname("   "))
	assert.Error(t, ValidateRealname("Bob")) // Assuming ValidateFL requires First and Last
	assert.NoError(t, ValidateRealname("Bob Dylan"))
}

func TestUpdateUserView(t *testing.T) {
	origRunForm := runForm
	defer func() { runForm = origRunForm }()

	runForm = func(f *huh.Form) error {
		return nil // mock successful form execution
	}

	view := &UpdateView{}
	err := UpdateUserView(view)
	assert.NoError(t, err)
}

func TestUpdateUserView_Error(t *testing.T) {
	origRunForm := runForm
	defer func() { runForm = origRunForm }()

	runForm = func(f *huh.Form) error {
		return errors.New("form error") // mock form error
	}

	view := &UpdateView{}
	err := UpdateUserView(view)
	assert.Error(t, err)
	assert.Equal(t, "form error", err.Error())
}
