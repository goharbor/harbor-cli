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

package project

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project/member"
	"github.com/spf13/cobra"
)

func Member() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member",
		Short:   `Manage members in a Project`,
		Long:    "Manage members and assign roles to them",
		Example: `  harbor member list`,
	}
	cmd.AddCommand(
		member.ListMemberCommand(),
		member.CreateMemberCommand(),
		member.DeleteMemberCommand(),
		member.UpdateMemberCommand(),
	)

	return cmd
}
