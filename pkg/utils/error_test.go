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

package utils

import (
	"errors"
	"testing"
)

type mockPayload struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

type mockErrorWithPayload struct {
	Payload mockPayload
	msg     string
}

func (e mockErrorWithPayload) Error() string {
	return e.msg
}

func TestParseHarborErrorMsg(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "generic error",
			err:      errors.New("generic error"),
			expected: "generic error",
		},
		{
			name: "error with payload",
			err: mockErrorWithPayload{
				Payload: mockPayload{
					Errors: []struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						{
							Code:    "404",
							Message: "not found",
						},
					},
				},
				msg: "some error message",
			},
			expected: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseHarborErrorMsg(tt.err)
			if actual != tt.expected {
				t.Errorf("ParseHarborErrorMsg() = %v, expected %v", actual, tt.expected)
			}
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
			name:     "error with brackets format",
			err:      errors.New("[GET /api/v2.0/projects][404] Project not found"),
			expected: "404",
		},
		{
			name:     "error with status format",
			err:      errors.New("failed to call api (status 500)"),
			expected: "500",
		},
		{
			name:     "generic error",
			err:      errors.New("connection refused"),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseHarborErrorCode(tt.err)
			if actual != tt.expected {
				t.Errorf("ParseHarborErrorCode() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}
