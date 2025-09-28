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
package root

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/password/change"
	"github.com/spf13/cobra"
)

func PasswordCommand() *cobra.Command {
	var opts change.PasswordChangeView

	cmd := &cobra.Command{
		Use:   "password",
		Short: "Change your password",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			change.ChangePasswordView(&opts)

			err := UpdatePassword(&opts)
			if err != nil {
				return fmt.Errorf("error changing password:%v", err)
			}
			return nil
		},
	}

	return cmd
}

func UpdatePassword(opts *change.PasswordChangeView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	userResp, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
	if err != nil {
		return err
	}

	_, err = client.User.UpdateUserPassword(ctx, &user.UpdateUserPasswordParams{
		Password: &models.PasswordReq{
			OldPassword: opts.OldPassword,
			NewPassword: opts.NewPassword,
		},
		UserID: userResp.Payload.UserID,
	})
	if err != nil {
		return err
	}
	if err := utils.GenerateEncryptionKey(); err != nil {
		fmt.Println("Encryption key already exists or could not be created:", err)
	}

	key, err := utils.GetEncryptionKey()
	if err != nil {
		return fmt.Errorf("failed to get encryption key: %s", err)
	}

	encryptedPassword, err := utils.Encrypt(key, []byte(opts.NewPassword))
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %s", err)
	}

	config, err := utils.GetCurrentHarborConfig()
	if err != nil {
		return fmt.Errorf("failed to get current Harbor config: %v", err)
	}

	credentialName := config.CurrentCredentialName

	existingCred, _ := utils.GetCredentials(credentialName)
	cred := utils.Credential{
		Name:          existingCred.Name,
		Username:      existingCred.Username,
		Password:      encryptedPassword,
		ServerAddress: existingCred.ServerAddress,
	}
	harborData, err := utils.GetCurrentHarborData()
	if err != nil {
		return fmt.Errorf("failed to get current harbor data: %s", err)
	}
	configPath := harborData.ConfigPath

	if err = utils.UpdateCredentialsInConfigFile(cred, configPath); err != nil {
		return fmt.Errorf("failed to update credentials: %s", err)
	}
	return nil
}
