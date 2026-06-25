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

package policy

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/preheat"
	"github.com/goharbor/harbor-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePolicyCommand_NilPolicy(t *testing.T) {
	originalGetPreheatPolicy := getPreheatPolicyFunc
	t.Cleanup(func() {
		getPreheatPolicyFunc = originalGetPreheatPolicy
	})

	getPreheatPolicyFunc = func(projectName, policyName string) (*preheat.GetPolicyOK, error) {
		return nil, nil
	}

	err := testutil.TestCmd(t, UpdatePolicyCommand, "my-project", "my-policy")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payload is empty")
}

func TestUpdatePolicyCommand_NilPayload(t *testing.T) {
	originalGetPreheatPolicy := getPreheatPolicyFunc
	t.Cleanup(func() {
		getPreheatPolicyFunc = originalGetPreheatPolicy
	})

	getPreheatPolicyFunc = func(projectName, policyName string) (*preheat.GetPolicyOK, error) {
		return &preheat.GetPolicyOK{Payload: nil}, nil
	}

	err := testutil.TestCmd(t, UpdatePolicyCommand, "my-project", "my-policy")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payload is empty")
}
