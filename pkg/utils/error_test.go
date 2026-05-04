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
	"fmt"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type mockSwaggerError struct {
	Payload interface{}
	message string
}

func (e *mockSwaggerError) Error() string {
	return e.message
}

func TestParseHarborErrorMsg_StructuredPayload(t *testing.T) {
	payload := struct {
		Errors []struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}{
		Errors: []struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			{Code: "NOT_FOUND", Message: "project not found"},
		},
	}

	err := &mockSwaggerError{
		Payload: payload,
		message: "[GET /projects/test][404] not found",
	}

	result := utils.ParseHarborErrorMsg(err)
	assert.Equal(t, "project not found", result)
}

func TestParseHarborErrorMsg_EmptyPayload(t *testing.T) {
	err := &mockSwaggerError{
		Payload: struct{}{},
		message: "response status code does not match any response statuses defined for this endpoint in the swagger spec (status 502): {}",
	}
	result := utils.ParseHarborErrorMsg(err)
	assert.Contains(t, result, "server error")
	assert.Contains(t, result, "HTTP 502")
}

func TestParseHarborErrorMsg_5xxErrors(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "502 Bad Gateway",
			errMsg:   "[GET /api/v2.0/users/current][502] response status code does not match any response statuses defined for this endpoint in the swagger spec (status 502): {}",
			expected: "server error (HTTP 502)",
		},
		{
			name:     "503 Service Unavailable",
			errMsg:   "[GET /projects][503] service unavailable",
			expected: "server error (HTTP 503)",
		},
		{
			name:     "500 Internal Server Error",
			errMsg:   "[POST /projects][500] internal server error",
			expected: "server error (HTTP 500)",
		},
		{
			name:     "5xx from status text pattern",
			errMsg:   "response status code does not match any response statuses defined for this endpoint in the swagger spec (status 504): {}",
			expected: "server error (HTTP 504)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errMsg)
			result := utils.ParseHarborErrorMsg(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHarborErrorMsg_4xxErrors(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "401 Unauthorized",
			errMsg:   "[GET /users/current][401] unauthorized",
			expected: "credentials invalid (HTTP 401)",
		},
		{
			name:     "403 Forbidden",
			errMsg:   "[GET /projects][403] forbidden",
			expected: "credentials invalid (HTTP 403)",
		},
		{
			name:     "404 Not Found",
			errMsg:   "[GET /projects/test][404] not found",
			expected: "request failed (HTTP 404)",
		},
		{
			name:     "409 Conflict",
			errMsg:   "[POST /projects][409] conflict",
			expected: "request failed (HTTP 409)",
		},
		{
			name:     "4xx from status text pattern",
			errMsg:   "response status code does not match any response statuses defined for this endpoint in the swagger spec (status 400): {}",
			expected: "request failed (HTTP 400)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errMsg)
			result := utils.ParseHarborErrorMsg(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHarborErrorMsg_UnrecognizedStatus(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "302 Found",
			errMsg:   "[GET /][302] found",
			expected: "unexpected HTTP 302 response",
		},
		{
			name:     "201 Created",
			errMsg:   "[POST /projects][201] created",
			expected: "unexpected HTTP 201 response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errMsg)
			result := utils.ParseHarborErrorMsg(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHarborErrorMsg_NoStatusCode(t *testing.T) {
	err := fmt.Errorf("some generic error without status code")
	result := utils.ParseHarborErrorMsg(err)
	assert.Equal(t, "some generic error without status code", result)
}

func TestParseHarborErrorMsg_NilError(t *testing.T) {
	result := utils.ParseHarborErrorMsg(nil)
	assert.Empty(t, result)
}

func TestParseHarborErrorCode_BracketedFormat(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "GET with 404",
			errMsg:   "[GET /projects/test][404] not found",
			expected: "404",
		},
		{
			name:     "POST with 409",
			errMsg:   "[POST /projects][409] conflict",
			expected: "409",
		},
		{
			name:     "DELETE with 200",
			errMsg:   "[DELETE /projects/test][200] ok",
			expected: "200",
		},
		{
			name:     "PUT with 500",
			errMsg:   "[PUT /projects][500] internal error",
			expected: "500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errMsg)
			result := utils.ParseHarborErrorCode(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHarborErrorCode_StatusTextFormat(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "status 502",
			errMsg:   "response status code does not match any response statuses defined for this endpoint in the swagger spec (status 502): {}",
			expected: "502",
		},
		{
			name:     "status 404",
			errMsg:   "something went wrong (status 404)",
			expected: "404",
		},
		{
			name:     "status 200",
			errMsg:   "unexpected (status 200) response",
			expected: "200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errMsg)
			result := utils.ParseHarborErrorCode(err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHarborErrorCode_NoMatch(t *testing.T) {
	err := errors.New("this is just a regular error message")
	result := utils.ParseHarborErrorCode(err)
	assert.Empty(t, result)
}
