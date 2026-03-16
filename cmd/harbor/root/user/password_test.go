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
package user

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestUserPasswordChangeCmd_Metadata(t *testing.T) {
	cmd := UserPasswordChangeCmd()

	if cmd == nil {
		t.Fatal("command should not be nil")
	}

	if cmd.Use != "password" {
		t.Fatalf("expected command 'password', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Fatal("Short description should not be empty")
	}
}

func TestUserPasswordChangeCmd_RunExists(t *testing.T) {
	cmd := UserPasswordChangeCmd()

	if cmd.Run == nil {
		t.Fatal("Run function should be defined")
	}
}

func TestUserPasswordChangeCmd_IsCobraCommand(t *testing.T) {
	cmd := UserPasswordChangeCmd()

	if _, ok := interface{}(cmd).(*cobra.Command); !ok {
		t.Fatal("expected cobra command")
	}
}
