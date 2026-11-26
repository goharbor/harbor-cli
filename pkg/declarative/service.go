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

package declarative

import (
	"context"
	"fmt"
)

// Backend reads and mutates Harbor resources for declarative reconciliation.
type Backend interface {
	Snapshot(context.Context) (*Configuration, error)
	Apply(context.Context, *Configuration, Action) error
}

// Service exports, plans, and applies declarative Harbor configuration.
type Service struct {
	backend Backend
}

// NewService creates a declarative configuration service.
func NewService(backend Backend) *Service {
	return &Service{backend: backend}
}

// Export returns normalized current state.
func (s *Service) Export(ctx context.Context) (*Configuration, error) {
	configuration, err := s.backend.Snapshot(ctx)
	if err != nil {
		return nil, fmt.Errorf("export Harbor configuration: %w", err)
	}
	configuration.Normalize()
	if err := configuration.Validate(); err != nil {
		return nil, fmt.Errorf("validate exported configuration: %w", err)
	}
	return configuration, nil
}

// Plan compares a desired document with current Harbor state.
func (s *Service) Plan(ctx context.Context, desired *Configuration) (Plan, error) {
	current, err := s.backend.Snapshot(ctx)
	if err != nil {
		return Plan{}, fmt.Errorf("read current Harbor configuration: %w", err)
	}
	plan, err := BuildPlan(desired, current)
	if err != nil {
		return Plan{}, fmt.Errorf("build reconciliation plan: %w", err)
	}
	return plan, nil
}

// Apply executes the mutating steps of an existing plan in dependency order.
func (s *Service) Apply(ctx context.Context, desired *Configuration, plan Plan) error {
	for _, action := range plan.Actions {
		if action.Operation == OperationNoop {
			continue
		}
		if err := s.backend.Apply(ctx, desired, action); err != nil {
			return fmt.Errorf("%s: %w", action, err)
		}
	}
	return nil
}
