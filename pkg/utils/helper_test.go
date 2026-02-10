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
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestFormatCreatedTime(t *testing.T) {
	// "now" should be 0 minutes ago
	now := time.Now().Format(time.RFC3339Nano)
	s, err := utils.FormatCreatedTime(now)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, "minute ago") || strings.HasSuffix(s, "minutes ago"))

	// invalid timestamp returns error
	_, err = utils.FormatCreatedTime("not-a-timestamp")
	assert.Error(t, err)
}

func TestFormatUrl(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"https://example.com", "https://example.com"},
		{"http://foo", "http://foo"},
		{"https://bar", "https://bar"},
		{"demo.goharbor.io", "https://demo.goharbor.io"},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got := utils.FormatUrl(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFormatSize(t *testing.T) {
	// 1048576 bytes == 1 MiB
	assert.Equal(t, "1.00MiB", utils.FormatSize(1024*1024))
	assert.Equal(t, "0B", utils.FormatSize(0))
}

func TestValidateUserName(t *testing.T) {
	assert.True(t, utils.ValidateUserName("alice"))
	assert.False(t, utils.ValidateUserName(""))
	assert.False(t, utils.ValidateUserName(strings.Repeat("x", 300)))
	assert.False(t, utils.ValidateUserName(`bad"name`))
}

func TestValidateEmail(t *testing.T) {
	assert.True(t, utils.ValidateEmail("foo@bar.com"))
	assert.False(t, utils.ValidateEmail("not-an-email"))
}

func TestValidateConfigPath(t *testing.T) {
	assert.True(t, utils.ValidateConfigPath("foo.yaml"))
	assert.True(t, utils.ValidateConfigPath("path/to/config.yml"))
	assert.False(t, utils.ValidateConfigPath("noext"))
	assert.False(t, utils.ValidateConfigPath("bad!.yaml"))
}

func TestValidateFL(t *testing.T) {
	assert.True(t, utils.ValidateFL("John Doe"))
	assert.False(t, utils.ValidateFL("SingleName"))
	assert.True(t, utils.ValidateFL("LongFirstname Lastname"))
}

func TestValidatePassword(t *testing.T) {
	// too short
	err := utils.ValidatePassword("Ab1")
	assert.Error(t, err)

	// no lowercase
	err = utils.ValidatePassword("ABCDEF12")
	assert.Error(t, err)

	// no uppercase
	err = utils.ValidatePassword("abcdef12")
	assert.Error(t, err)

	// no digit
	err = utils.ValidatePassword("Abcdefgh")
	assert.Error(t, err)

	// valid
	err = utils.ValidatePassword("Abcd1234")
	assert.NoError(t, err)
}

func TestValidateTagName(t *testing.T) {
	assert.True(t, utils.ValidateTagName("v1.0.0"))
	assert.False(t, utils.ValidateTagName(".bad"))
}

func TestValidateProjectName(t *testing.T) {
	assert.True(t, utils.ValidateProjectName("project_1"))
	assert.False(t, utils.ValidateProjectName("-invalid"))
}

func TestValidateStorageLimit(t *testing.T) {
	assert.NoError(t, utils.ValidateStorageLimit("0"))
	assert.NoError(t, utils.ValidateStorageLimit("-1"))
	assert.Error(t, utils.ValidateStorageLimit("foo"))
	assert.Error(t, utils.ValidateStorageLimit("2048"))
}

func TestValidateRegistryName(t *testing.T) {
	assert.True(t, utils.ValidateRegistryName("registry01"))
	assert.False(t, utils.ValidateRegistryName("-bad"))
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// valid URLs
		{"valid https domain", "https://demo.goharbor.io", false},
		{"valid http domain", "http://registry.example.com", false},
		{"valid domain with port", "https://harbor.local:8443", false},
		{"valid domain with path", "https://example.com/v2", false},
		{"valid IPv4", "http://192.168.1.1", false},
		{"valid IPv4 with port", "http://192.168.1.1:8080", false},
		{"valid IPv6", "http://[::1]", false},

		// invalid URLs
		{"empty string", "", true},
		{"no host", "https://", true},
		{"bare path", "/just/a/path", true},
		{"invalid domain double dot", "https://invalid..domain", true},
		{"domain starting with hyphen", "https://-invalid.com", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.ValidateURL(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPrintFormat(t *testing.T) {
	// capture stdout
	var buf bytes.Buffer
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	type payload struct {
		Foo string `json:"foo"`
	}
	obj := payload{Foo: "bar"}

	// JSON
	err := utils.PrintFormat(obj, "json")
	w.Close()
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to capture output: %v", err)
	}
	os.Stdout = oldOut
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"foo": "bar"`)

	// YAML
	buf.Reset()
	r, w, _ = os.Pipe()
	os.Stdout = w
	err = utils.PrintFormat(obj, "yaml")
	w.Close()
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to capture output: %v", err)
	}
	os.Stdout = oldOut
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "foo: bar")

	// unsupported
	err = utils.PrintFormat(obj, "xml")
	assert.Error(t, err)
}
