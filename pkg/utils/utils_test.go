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
	"fmt"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func Test_Sanitize_ServerAddress(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://harbor.example.com", "https-harbor-example-com"},
		{"http://harbor.example.com", "http-harbor-example-com"},
		{"https://harbor.example.com:8080", "https-harbor-example-com-8080"},
		{"http://harbor.example.com:8080", "http-harbor-example-com-8080"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := utils.SanitizeServerAddress(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func Test_ParseProjectRepoReference(t *testing.T) {
	tests := []struct {
		input               string
		expProject, expRepo string
		expReference        string
		wantErr             bool
	}{
		{"project/repo:reference", "project", "repo", "reference", false},
		{"project/repo:reference:tag", "project", "repo:reference", "tag", false},
		{"project/repo:reference@sha256:1234567890abcdef", "project", "repo:reference", "sha256:1234567890abcdef", false},
		{"project/repo", "project", "repo", "", true},
		{"project/repo/reference", "project", "repo/reference", "", true},
		{"bad-format", "", "", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			project, repo, reference, err := utils.ParseProjectRepoReference(tc.input)
			fmt.Println(project, repo, reference, err)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			// no error expected
			assert.NoError(t, err)
			assert.Equal(t, tc.expProject, project)
			assert.Equal(t, tc.expRepo, repo)
			assert.Equal(t, tc.expReference, reference)
		})
	}
}

func Test_ParseProjectRepo(t *testing.T) {
	tests := []struct {
		input               string
		expProject, expRepo string
		wantErr             bool
	}{
		{"project/repo/reference", "project", "repo/reference", false},
		{"project/repo", "project", "repo", false},
		{"bad-format", "", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			project, repo, err := utils.ParseProjectRepo(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			// no error expected
			assert.NoError(t, err)
			assert.Equal(t, tc.expProject, project)
			assert.Equal(t, tc.expRepo, repo)
		})
	}
}

func TestStorageStringToBytes(t *testing.T) {
	// Valid inputs
	tests := []struct {
		input    string
		expected int64
	}{
		{"1MiB", 1024 * 1024},
		{"1GiB", 1024 * 1024 * 1024},
		{"1TiB", 1024 * 1024 * 1024 * 1024},
	}

	for _, test := range tests {
		result, err := utils.StorageStringToBytes(test.input)
		assert.NoError(t, err, "Unexpected error for input %s", test.input)
		assert.Equal(
			t,
			test.expected,
			result,
			"Expected %d but got %d for input %s",
			test.expected,
			result,
			test.input,
		)
	}

	// Invalid inputs
	invalidInputs := []string{
		"1KB",
		"1000",
		"10PB",
		"1GiBGiB",
		"1.03GiB",
		"1.08TiB",
	}

	for _, input := range invalidInputs {
		_, err := utils.StorageStringToBytes(input)
		assert.Error(t, err, "Expected error for input %s but got none", input)
	}

	// Exceeding maximum value
	_, err := utils.StorageStringToBytes("1025TiB")
	assert.Error(t, err, "Expected error for input exceeding 1024TiB but got none")
}
