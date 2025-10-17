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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func TestDefaultCredentialName(t *testing.T) {
	name := utils.DefaultCredentialName("john", "https://harbor.example.com:8080")
	assert.Equal(t, "john@https-harbor-example-com-8080", name)
}

func TestToKebabCase(t *testing.T) {
	assert.Equal(t, "hello-world", utils.ToKebabCase("Hello World"))
	assert.Equal(t, "multi-space-string", utils.ToKebabCase("Multi Space String"))
	assert.Equal(t, "already-kebab", utils.ToKebabCase("already-kebab"))
}

func TestCapitalize(t *testing.T) {
	// current behavior returns the string unchanged, ensure we lock that in
	assert.Equal(t, "", utils.Capitalize(""))
	assert.Equal(t, "harbor", utils.Capitalize("harbor"))
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	f()
	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPrintPayloadInJSONFormat(t *testing.T) {
	type payloadT struct {
		Name string `json:"name" yaml:"name"`
		ID   int    `json:"id" yaml:"id"`
	}
	p := payloadT{Name: "harbor", ID: 1}

	expectedBytes, _ := json.MarshalIndent(p, "", "  ")
	expected := string(expectedBytes) + "\n" // Println adds a newline
	out := captureStdout(func() { utils.PrintPayloadInJSONFormat(p) })
	assert.Equal(t, expected, out)

	// nil should print nothing
	out = captureStdout(func() { utils.PrintPayloadInJSONFormat(nil) })
	assert.Equal(t, "", out)
}

func TestPrintPayloadInYAMLFormat(t *testing.T) {
	type payloadT struct {
		Name string `json:"name" yaml:"name"`
		ID   int    `json:"id" yaml:"id"`
	}
	p := payloadT{Name: "harbor", ID: 1}

	// YAML marshaling is deterministic for a struct; replicate function behavior
	out := captureStdout(func() { utils.PrintPayloadInYAMLFormat(p) })
	assert.Contains(t, out, "name: harbor")
	assert.Contains(t, out, "id: 1")

	// nil should print nothing
	out = captureStdout(func() { utils.PrintPayloadInYAMLFormat(nil) })
	assert.Equal(t, "", out)
}

func TestSavePayloadJSON(t *testing.T) {
	type payloadT struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
	p := payloadT{Name: "harbor", ID: 1}

	dir := t.TempDir()
	base := filepath.Join(dir, "out")
	utils.SavePayloadJSON(base, p)

	path := base + ".json"
	b, err := os.ReadFile(path)
	assert.NoError(t, err)

	expected, _ := json.MarshalIndent(p, "", "  ")
	assert.Equal(t, string(expected), string(b))

	fi, err := os.Stat(path)
	assert.NoError(t, err)
	// File mode should be 0600
	assert.Equal(t, os.FileMode(0o600), fi.Mode().Perm())
}
