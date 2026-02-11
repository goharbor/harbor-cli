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
	"fmt"
	"strings"
)

// BuildSortParam validates sort fields against allowed fields and returns
// a comma-separated sort string for the API.
// Each sort field can optionally be prefixed with '-' for descending order.
func BuildSortParam(sortFields []string, validFields []string) (string, error) {
	for _, field := range sortFields {
		// Strip leading '-' (descending indicator) for validation
		fieldName := strings.TrimPrefix(field, "-")
		if err := validateKey(fieldName, validFields); err != nil {
			return "", fmt.Errorf("invalid sort field: %s, supported fields are: %s", fieldName, strings.Join(validFields, ", "))
		}
	}

	return strings.Join(sortFields, ","), nil
}
