package schedules

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

import (
	"errors"
	"strings"
	"testing"
)

func TestSchedulesCommandIncludesExpectedSubcommands(t *testing.T) {
	cmd := SchedulesCommand()

	got := make(map[string]struct{}, len(cmd.Commands()))
	for _, subcommand := range cmd.Commands() {
		got[subcommand.Name()] = struct{}{}
	}

	for _, want := range []string{"list", "status", "pause-all", "resume-all"} {
		if _, ok := got[want]; !ok {
			t.Fatalf("expected subcommand %q to be registered", want)
		}
	}
}

func TestListCommandDefaults(t *testing.T) {
	cmd := ListCommand()

	if got := cmd.Flags().Lookup("page").DefValue; got != "1" {
		t.Fatalf("expected default page to be 1, got %s", got)
	}

	if got := cmd.Flags().Lookup("page-size").DefValue; got != "20" {
		t.Fatalf("expected default page-size to be 20, got %s", got)
	}
}

func TestFormatScheduleError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		permission string
		wantSubstr string
	}{
		{
			name:       "bad request",
			err:        errors.New("[GET /schedule][400] (status 400)"),
			permission: "ActionStop",
			wantSubstr: "invalid request",
		},
		{
			name:       "authentication required",
			err:        errors.New("[GET /schedule][401] (status 401)"),
			permission: "ActionStop",
			wantSubstr: "authentication required",
		},
		{
			name:       "authenticated but forbidden",
			err:        errors.New("[GET /schedule][403] (status 403)"),
			permission: "authenticated",
			wantSubstr: "authenticated but lacks access",
		},
		{
			name:       "fallback message",
			err:        errors.New("unexpected failure"),
			permission: "ActionStop",
			wantSubstr: "unexpected failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatScheduleError("failed to test", tt.err, tt.permission)
			if got == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(got.Error(), tt.wantSubstr) {
				t.Fatalf("expected %q to contain %q", got.Error(), tt.wantSubstr)
			}
		})
	}
}
