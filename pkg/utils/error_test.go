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
	"errors"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type mockPayload struct {
	Payload interface{}
}

func (m *mockPayload) Error() string {
	return "mock error"
}

func TestParseHarborErrorMsg(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			"nil error",
			nil,
			"",
		},
		{
			"plain error",
			errors.New("simple error"),
			"simple error",
		},
		{
			"single harbor error",
			&mockPayload{
				Payload: utils.HarborErrorPayload{
					Errors: []struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						{Code: "400", Message: "first error"},
					},
				},
			},
			"first error",
		},
		{
			"multiple harbor errors",
			&mockPayload{
				Payload: utils.HarborErrorPayload{
					Errors: []struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						{Code: "400", Message: "first error"},
						{Code: "401", Message: "second error"},
					},
				},
			},
			"first error; second error",
		},
		{
			"invalid payload type",
			&mockPayload{
				Payload: make(chan int),
			},
			"mock error",
		},
		{
			"invalid json format",
			&mockPayload{
				Payload: "this is a string, not a struct",
			},
			"mock error",
		},
		{
			"empty errors slice",
			&mockPayload{
				Payload: utils.HarborErrorPayload{
					Errors: nil,
				},
			},
			"mock error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.ParseHarborErrorMsg(tt.err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParseHarborErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			"bracket format",
			errors.New("[GET /projects][404] Not Found"),
			"404",
		},
		{
			"status format",
			errors.New("request failed (status 401) Unauthorized"),
			"401",
		},
		{
			"no code",
			errors.New("generic error"),
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.ParseHarborErrorCode(tt.err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
