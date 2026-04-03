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

package jobs

import "testing"

func TestJobsCommandIncludesLogSubcommand(t *testing.T) {
	cmd := JobsCommand()

	if cmd.Use != "jobs" {
		t.Fatalf("expected command use to be jobs, got %s", cmd.Use)
	}

	if len(cmd.Commands()) != 1 {
		t.Fatalf("expected one subcommand, got %d", len(cmd.Commands()))
	}

	if got := cmd.Commands()[0].Name(); got != "log" {
		t.Fatalf("expected log subcommand, got %s", got)
	}
}

func TestLogCommandRequiresJobID(t *testing.T) {
	cmd := LogCommand()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil || err.Error() != "--job-id must be specified" {
		t.Fatalf("expected job-id validation error, got %v", err)
	}
}