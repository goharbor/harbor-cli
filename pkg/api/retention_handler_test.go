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

package api

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/retention/create"
)

func TestGenerateRetentionRule(t *testing.T) {
	tests := []struct {
		name        string
		opts        create.CreateView
		expectError bool
		checkRule   func(*models.RetentionRule) error
	}{
		{
			name: "always template with no params",
			opts: create.CreateView{
				Action:    "retain",
				Algorithm: "or",
				Scope: create.RetentionPolicyScope{
					Level: "project",
					Ref:   1,
				},
				Template: "always",
				TagSelectors: create.RetentionSelector{
					Decoration: "matches",
					Pattern:    "release-*",
				},
				ScopeSelectors: create.RetentionSelector{
					Decoration: "repoMatches",
					Pattern:    "**",
				},
			},
			expectError: false,
			checkRule: func(rule *models.RetentionRule) error {
				if rule.Action != "retain" {
					return &testError{msg: "action mismatch"}
				}
				if rule.Template != "always" {
					return &testError{msg: "template mismatch"}
				}
				if rule.Params != nil {
					return &testError{msg: "params should be nil for always template"}
				}
				return nil
			},
		},
		{
			name: "latestPushedK template with count param",
			opts: create.CreateView{
				Action:    "retain",
				Algorithm: "or",
				Scope: create.RetentionPolicyScope{
					Level: "project",
					Ref:   1,
				},
				Template: "latestPushedK",
				Params: create.ParamsValue{
					Name:  "count",
					Value: "5",
				},
				TagSelectors: create.RetentionSelector{
					Decoration: "matches",
					Pattern:    "release-*",
				},
				ScopeSelectors: create.RetentionSelector{
					Decoration: "repoMatches",
					Pattern:    "**",
				},
			},
			expectError: false,
			checkRule: func(rule *models.RetentionRule) error {
				if rule.Params == nil {
					return &testError{msg: "params should not be nil"}
				}
				if val, ok := rule.Params["latestPushedK"]; !ok || val != 5 {
					return &testError{msg: "param value mismatch"}
				}
				return nil
			},
		},
		{
			name: "nDaysSinceLastPush template with days param",
			opts: create.CreateView{
				Action:    "delete",
				Algorithm: "and",
				Scope: create.RetentionPolicyScope{
					Level: "repository",
					Ref:   2,
				},
				Template: "nDaysSinceLastPush",
				Params: create.ParamsValue{
					Name:  "days",
					Value: "30",
				},
				TagSelectors: create.RetentionSelector{
					Decoration: "matches",
					Pattern:    "*",
				},
				ScopeSelectors: create.RetentionSelector{
					Decoration: "repoMatches",
					Pattern:    "my-repo",
				},
			},
			expectError: false,
			checkRule: func(rule *models.RetentionRule) error {
				if rule.Action != "delete" {
					return &testError{msg: "action mismatch"}
				}
				if val, ok := rule.Params["nDaysSinceLastPush"]; !ok || val != 30 {
					return &testError{msg: "param value mismatch"}
				}
				return nil
			},
		},
		{
			name: "invalid param value for non-always template",
			opts: create.CreateView{
				Action:    "retain",
				Algorithm: "or",
				Scope: create.RetentionPolicyScope{
					Level: "project",
					Ref:   1,
				},
				Template: "latestPushedK",
				Params: create.ParamsValue{
					Name:  "count",
					Value: "invalid",
				},
				TagSelectors: create.RetentionSelector{
					Decoration: "matches",
					Pattern:    "release-*",
				},
				ScopeSelectors: create.RetentionSelector{
					Decoration: "repoMatches",
					Pattern:    "**",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := generateRetentionRule(tt.opts)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tt.checkRule != nil {
				if err := tt.checkRule(rule); err != nil {
					t.Errorf("rule check failed: %v", err)
				}
			}
		})
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}