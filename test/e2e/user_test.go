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

//go:build e2e

package e2e

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Harbor lists 10 users per page by default, so deleting more than 10 users
// by name checks that the lookup is not limited to the first page.
func Test_UserDelete_FindsUsersPastFirstPage(t *testing.T) {
	const userCount = 12

	loginAndReload(t)

	prefix := fmt.Sprintf("e2e-819-%d", time.Now().UnixNano())
	names := make([]string, userCount)
	ids := make([]int64, userCount)
	for i := range userCount {
		names[i] = fmt.Sprintf("%s-%02d", prefix, i)
		ids[i] = createUser(t, names[i])
	}

	require.NoError(t, harbor(t, append([]string{"user", "delete"}, names...)...))

	for i, id := range ids {
		require.False(t, userExists(t, id), "user %s (id %d) should have been deleted", names[i], id)
	}
}

func Test_UserDelete_UnknownUserFails(t *testing.T) {
	loginAndReload(t)

	name := fmt.Sprintf("e2e-819-missing-%d", time.Now().UnixNano())
	err := harbor(t, "user", "delete", name)
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("user %q not found", name))
}

func Test_UserDelete_ByNumericID(t *testing.T) {
	loginAndReload(t)

	name := fmt.Sprintf("e2e-819-byid-%d", time.Now().UnixNano())
	id := createUser(t, name)

	require.NoError(t, harbor(t, "user", "delete", strconv.FormatInt(id, 10)))
	require.False(t, userExists(t, id), "user %s (id %d) should have been deleted", name, id)
}

func Test_UserDelete_UnknownUserDeletesNothing(t *testing.T) {
	loginAndReload(t)

	name := fmt.Sprintf("e2e-819-keep-%d", time.Now().UnixNano())
	id := createUser(t, name)

	err := harbor(t, "user", "delete", name, name+"-missing")
	require.Error(t, err)
	require.True(t, userExists(t, id), "user %s (id %d) should not be deleted when another argument fails", name, id)
}

// Harbor reads a q value wrapped in brackets as a range, so such names need
// a literal lookup.
func Test_UserDelete_NameWithBrackets(t *testing.T) {
	loginAndReload(t)

	name := fmt.Sprintf("[e2e-819-%d]", time.Now().UnixNano())
	id := createUser(t, name)

	require.NoError(t, harbor(t, "user", "delete", name))
	require.False(t, userExists(t, id), "user %s (id %d) should have been deleted", name, id)
}

// Harbor URL-decodes the q filter a second time and treats a leading
// backslash, [, { or ( specially, so these names must still resolve exactly.
func Test_UserDelete_NamesWithQueryCharacters(t *testing.T) {
	loginAndReload(t)

	suffix := time.Now().UnixNano()
	names := []string{
		fmt.Sprintf("e2e+plus-%d", suffix),
		fmt.Sprintf(`\e2e-backslash-%d`, suffix),
		fmt.Sprintf("[e2e-bracket-%d]", suffix),
		fmt.Sprintf("e2e&amp-%d", suffix),
	}
	ids := make([]int64, len(names))
	for i, name := range names {
		ids[i] = createUser(t, name)
	}

	require.NoError(t, harbor(t, append([]string{"user", "delete"}, names...)...))

	for i, id := range ids {
		require.False(t, userExists(t, id), "user %s (id %d) should have been deleted", names[i], id)
	}
}
