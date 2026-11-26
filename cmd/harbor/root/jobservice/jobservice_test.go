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
package jobservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobService_WiresSubcommands(t *testing.T) {
	cmd := JobService()

	assert.Equal(t, "jobservice", cmd.Use)

	subcommands := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = true
	}

	assert.True(t, subcommands["queues"], "jobservice should wire queues")
	assert.True(t, subcommands["pools"], "jobservice should wire pools")
	assert.True(t, subcommands["workers"], "jobservice should wire workers")
}
