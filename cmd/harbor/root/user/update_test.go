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

func TestUserUpdateCmd_Metadata(t *testing.T) {
	cmd := UserUpdateCmd()

	if cmd == nil {
		t.Fatal("command should not be nil")
	}

	if cmd.Use != "update [USER_NAME_OR_ID]" {
		t.Fatalf("expected command 'update [USER_NAME_OR_ID]', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Fatal("Short description should not be empty")
	}
}

func TestUserUpdateCmd_RunExists(t *testing.T) {
	cmd := UserUpdateCmd()

	if cmd.RunE == nil {
		t.Fatal("Run function should be defined")
	}
}

func TestUserUpdateCmd_IsCobraCommand(t *testing.T) {
	cmd := UserUpdateCmd()

	if _, ok := interface{}(cmd).(*cobra.Command); !ok {
		t.Fatal("expected cobra command")
	}
}

func TestUserUpdateCmd_Flags(t *testing.T) {
	cmd := UserUpdateCmd()

	emailFlag := cmd.Flags().Lookup("email")
	if emailFlag == nil {
		t.Fatal("expected 'email' flag")
	}

	realnameFlag := cmd.Flags().Lookup("realname")
	if realnameFlag == nil {
		t.Fatal("expected 'realname' flag")
	}

	commentFlag := cmd.Flags().Lookup("comment")
	if commentFlag == nil {
		t.Fatal("expected 'comment' flag")
	}
}
