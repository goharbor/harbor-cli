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

// Builds the `q` param for List API's
func BuildQueryParam(fuzzy, match, ranges []string, validKeys []string) (string, error) {
	var parts []string

	// Fuzzy
	for _, v := range fuzzy {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			return "", fmt.Errorf("invalid fuzzy arg: %s ", v)
		}

		if err := validateKey(kv[0], validKeys); err != nil {
			return "", err
		}

		parts = append(parts, fmt.Sprintf("%s=~%s", kv[0], kv[1]))
	}

	// Exact Match's
	for _, v := range match {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			return "", fmt.Errorf("invalid match arg: %s ", v)
		}

		if err := validateKey(kv[0], validKeys); err != nil {
			return "", err
		}

		parts = append(parts, fmt.Sprintf("%s=%s", kv[0], kv[1]))
	}

	// Ranges
	for _, v := range ranges {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			return "", fmt.Errorf("invalid range arg: %s ", v)
		}

		if err := validateKey(kv[0], validKeys); err != nil {
			return "", err
		}

		// Validating that range is in format min~max
		rng := strings.Split(kv[1], "~")
		if len(rng) != 2 {
			return "", fmt.Errorf("invalid range arg: %s ", v)
		}

		parts = append(parts, fmt.Sprintf("%s=[%s~%s]", kv[0], rng[0], rng[1]))
	}

	return strings.Join(parts, ","), nil
}

// Validates Key provided by user for ListFlags.Q
func validateKey(key string, validKeys []string) error {
	found := false
	for _, v := range validKeys {
		if v == key {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("invalid key for query: %s, supported keys are: %s", key, strings.Join(validKeys, ", "))
	}

	return nil
}
