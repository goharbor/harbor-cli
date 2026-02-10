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
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type HarborErrorPayload struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func ParseHarborErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	val := reflect.ValueOf(err)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName("Payload")
	if field.IsValid() {
		payload := field.Interface()
		jsonBytes, jsonErr := json.Marshal(payload)
		if jsonErr == nil {
			var harborErr HarborErrorPayload
			if unmarshalErr := json.Unmarshal(jsonBytes, &harborErr); unmarshalErr == nil {
				if len(harborErr.Errors) > 0 {
					return harborErr.Errors[0].Message
				}
			}
		}
	}
	return fmt.Sprintf("%v", err.Error())
}

func ParseHarborErrorCode(err error) string {
	errStr := err.Error()

	// Try format: [METHOD /path][CODE] - e.g., [GET /projects][404]
	parts := strings.Split(errStr, "]")
	if len(parts) >= 2 {
		codePart := strings.TrimSpace(parts[1])
		if strings.HasPrefix(codePart, "[") && len(codePart) == 4 {
			code := codePart[1:4]
			return code
		}
	}

	// Try format: (status CODE) - e.g., (status 404)
	re := regexp.MustCompile(`\(status\s+(\d{3})\)`)
	if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
		return matches[1]
	}

	return ""
}
