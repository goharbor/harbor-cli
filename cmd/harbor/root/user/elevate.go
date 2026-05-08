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
	"errors"
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/spf13/cobra"
)

var (
	getUsersIDByName  = api.GetUsersIdByName
	getUserIDFromUser = prompt.GetUserIdFromUser
	confirmElevation  = views.ConfirmElevation
	elevateUserAPI    = api.ElevateUser
)

var elevateUseID bool

func ElevateUser(args []string) error {
	var err error
	var userId int64
	if len(args) > 0 {
		if elevateUseID {
			parsedID, parseErr := strconv.ParseInt(args[0], 10, 64)
			if parseErr != nil {
				return fmt.Errorf("invalid ID '%s': %v", args[0], parseErr)
			}
			userId = parsedID
		} else {
			userId, err = getUsersIDByName(args[0])
			if err != nil {
				return err
			}
			if userId == 0 {
				return fmt.Errorf("user '%s' not found", args[0])
			}
		}
	} else {
		userId, err = getUserIDFromUser()
		if err != nil {
			return fmt.Errorf("failed to get user id: %v", err)
		}
	}
	confirm, err := confirmElevation()
	if err != nil {
		return fmt.Errorf("failed to confirm elevation: %v", err)
	}
	if !confirm {
		return errors.New("user declined elevation")
	}

	err = elevateUserAPI(userId)
	if err != nil {
		if isUnauthorizedError(err) {
			return fmt.Errorf("permission denied: admin privileges are required: %w", err)
		}
		return fmt.Errorf("failed to elevate user: %v", err)
	}
	return nil
}

func ElevateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elevate",
		Short: "elevate user",
		Long:  "elevate user to admin role",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ElevateUser(args)
		},
	}

	cmd.Flags().BoolVar(&elevateUseID, "id", false, "Use ID instead of username")
	return cmd
}
