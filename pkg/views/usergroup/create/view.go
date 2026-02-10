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
package create

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

type CreateUserGroupInput struct {
	GroupName   string
	GroupType   int64
	LDAPGroupDN string
}

func CreateUserGroupView(opts *CreateUserGroupInput) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter new group name").
				Value(&opts.GroupName),
			huh.NewSelect[int64]().
				Title("Select group type").
				Options(
					huh.NewOption("LDAP", int64(1)),
					huh.NewOption("HTTP", int64(2)),
					huh.NewOption("OIDC", int64(3)),
				).
				Value(&opts.GroupType),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("form input error: %v", err)
	}

	// If LDAP
	if opts.GroupType == 1 {
		err := huh.NewInput().
			Title("LDAP Group DN").
			Value(&opts.LDAPGroupDN).
			Run()
		if err != nil {
			return fmt.Errorf("form input error: %v", err)
		}
	}

	return nil
}
