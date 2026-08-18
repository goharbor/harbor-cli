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

package declarative_test

import (
	"context"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/declarative"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordingBackend struct {
	snapshot *declarative.Configuration
	actions  []declarative.Action
}

func (b *recordingBackend) Snapshot(context.Context) (*declarative.Configuration, error) {
	return b.snapshot, nil
}

func (b *recordingBackend) Apply(_ context.Context, _ *declarative.Configuration, action declarative.Action) error {
	b.actions = append(b.actions, action)
	return nil
}

func TestServiceApplySkipsNoopActions(t *testing.T) {
	backend := &recordingBackend{snapshot: declarative.NewConfiguration()}
	service := declarative.NewService(backend)
	desired := declarative.NewConfiguration()
	plan := declarative.Plan{Actions: []declarative.Action{
		{Operation: declarative.OperationNoop, Resource: "project", Name: "existing"},
		{Operation: declarative.OperationCreate, Resource: "project", Name: "new"},
	}}

	require.NoError(t, service.Apply(context.Background(), desired, plan))
	require.Len(t, backend.actions, 1)
	assert.Equal(t, "new", backend.actions[0].Name)
}
