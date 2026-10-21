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

package declarative_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/declarative"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAMLRoundTrip(t *testing.T) {
	public := true
	document := declarative.NewConfiguration()
	document.Spec.Projects = []declarative.Project{{Name: "zeta"}, {Name: "alpha", Public: &public}}

	var output bytes.Buffer
	require.NoError(t, declarative.Encode(&output, document, declarative.FormatYAML))

	decoded, err := declarative.Decode(strings.NewReader(output.String()), declarative.FormatYAML)
	require.NoError(t, err)
	assert.Equal(t, "alpha", decoded.Spec.Projects[0].Name)
	assert.Contains(t, output.String(), "apiVersion: goharbor.io/v1alpha1")
}

func TestDecodeRejectsUnknownFields(t *testing.T) {
	document := `apiVersion: goharbor.io/v1alpha1
kind: HarborConfiguration
spec:
  unknown: true
`

	_, err := declarative.Decode(strings.NewReader(document), declarative.FormatYAML)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "field unknown not found")
}

func TestDecodeRejectsMultipleDocuments(t *testing.T) {
	document := `apiVersion: goharbor.io/v1alpha1
kind: HarborConfiguration
---
apiVersion: goharbor.io/v1alpha1
kind: HarborConfiguration
`

	_, err := declarative.Decode(strings.NewReader(document), declarative.FormatYAML)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "multiple documents")
}

func TestValidateRejectsDuplicateNames(t *testing.T) {
	document := declarative.NewConfiguration()
	document.Spec.Registries = []declarative.Registry{{Name: "docker-hub"}, {Name: "docker-hub"}}

	err := document.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), `duplicate registry "docker-hub"`)
}

func TestExampleConfigurationStaysValid(t *testing.T) {
	filename := filepath.Join("..", "..", "examples", "declarative", "harbor.yaml")

	document, err := declarative.ReadFile(filename)
	require.NoError(t, err)
	assert.Equal(t, declarative.APIVersion, document.APIVersion)
	assert.NotEmpty(t, document.Spec.Projects)
}
