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
package instance

import (
	"fmt"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	viewupdate "github.com/goharbor/harbor-cli/pkg/views/instance/update"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	Name         string
	Description  string
	Endpoint     string
	AuthMode     string
	Enabled      bool
	Insecure     bool
	AuthUsername string
	AuthPassword string
	AuthToken    string
}

func UpdateInstanceCommand() *cobra.Command {
	var opts updateOptions
	var isID bool

	cmd := &cobra.Command{
		Use:   "update [NAME|ID]",
		Short: "Update a preheat provider instance in Harbor",
		Long: `Update a preheat provider instance in Harbor by name or ID. If no update
flags are provided, the command opens an interactive update form.`,
		Example: `  harbor-cli instance update my-instance --description "Updated preheat instance"
  harbor-cli instance update 1 --id --enable=false
  harbor-cli instance update my-instance --authmode BASIC --auth-username admin --auth-password Harbor12345`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var instanceName string

			if isID && len(args) == 0 {
				return fmt.Errorf("instance ID must be provided when using --id")
			}

			if len(args) > 0 {
				instanceName = args[0]
			} else {
				instanceName, err = prompt.GetInstanceNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get instance name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			resp, err := api.GetInstance(instanceName, isID)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("instance %s does not exist", instanceName)
				}
				return fmt.Errorf("failed to get instance: %v", utils.ParseHarborErrorMsg(err))
			}
			if resp == nil || resp.Payload == nil {
				return fmt.Errorf("failed to get instance: empty response")
			}

			instance := resp.Payload
			originalName := instance.Name

			if hasUpdateFlagChanges(cmd) {
				err = applyUpdateFlags(cmd, instance, opts)
			} else {
				err = viewupdate.UpdateInstanceView(instance)
			}
			if err != nil {
				return err
			}

			err = api.UpdateInstance(originalName, *instance)
			if err != nil {
				return fmt.Errorf("failed to update instance: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Printf("Instance '%s' updated successfully\n", instance.Name)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Get instance by id")
	flags.StringVarP(&opts.Name, "name", "n", "", "New name for the instance")
	flags.StringVarP(&opts.Endpoint, "url", "u", "", "Endpoint URL for the instance")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the instance")
	flags.BoolVarP(&opts.Insecure, "insecure", "i", false, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.BoolVarP(&opts.Enabled, "enable", "", false, "Whether the instance is enabled or not")
	flags.StringVarP(&opts.AuthMode, "authmode", "a", "", "Authentication mode (NONE, BASIC, OAUTH)")
	flags.StringVar(&opts.AuthUsername, "auth-username", "", "Username for BASIC authentication")
	flags.StringVar(&opts.AuthPassword, "auth-password", "", "Password for BASIC authentication")
	flags.StringVar(&opts.AuthToken, "auth-token", "", "Token for OAUTH authentication")

	return cmd
}

func hasUpdateFlagChanges(cmd *cobra.Command) bool {
	flags := cmd.Flags()
	return flags.Changed("name") ||
		flags.Changed("url") ||
		flags.Changed("description") ||
		flags.Changed("insecure") ||
		flags.Changed("enable") ||
		flags.Changed("authmode") ||
		flags.Changed("auth-username") ||
		flags.Changed("auth-password") ||
		flags.Changed("auth-token")
}

func applyUpdateFlags(cmd *cobra.Command, instance *models.Instance, opts updateOptions) error {
	flags := cmd.Flags()

	if flags.Changed("name") {
		if strings.TrimSpace(opts.Name) == "" {
			return fmt.Errorf("name cannot be empty or only spaces")
		}
		instance.Name = strings.TrimSpace(opts.Name)
	}
	if flags.Changed("url") {
		formattedURL := utils.FormatUrl(opts.Endpoint)
		if err := utils.ValidateURL(formattedURL); err != nil {
			return err
		}
		instance.Endpoint = formattedURL
	}
	if flags.Changed("description") {
		instance.Description = opts.Description
	}
	if flags.Changed("insecure") {
		instance.Insecure = opts.Insecure
	}
	if flags.Changed("enable") {
		instance.Enabled = opts.Enabled
	}

	if flags.Changed("auth-username") || flags.Changed("auth-password") || flags.Changed("auth-token") {
		if !flags.Changed("authmode") {
			return fmt.Errorf("authmode is required when updating auth credentials")
		}
	}

	if !flags.Changed("authmode") {
		return nil
	}

	instance.AuthMode = strings.ToUpper(strings.TrimSpace(opts.AuthMode))

	switch instance.AuthMode {
	case "BASIC":
		if strings.TrimSpace(opts.AuthUsername) == "" || strings.TrimSpace(opts.AuthPassword) == "" {
			return fmt.Errorf("username and password are required when authmode is BASIC. Use --auth-username and --auth-password flags")
		}
		instance.AuthInfo = map[string]string{
			"username": strings.TrimSpace(opts.AuthUsername),
			"password": opts.AuthPassword,
		}
	case "OAUTH":
		if strings.TrimSpace(opts.AuthToken) == "" {
			return fmt.Errorf("token is required when authmode is OAUTH. Use --auth-token flag")
		}
		instance.AuthInfo = map[string]string{
			"token": strings.TrimSpace(opts.AuthToken),
		}
	case "NONE":
		instance.AuthInfo = nil
	default:
		return fmt.Errorf("invalid authmode '%s'. Valid options: NONE, BASIC, OAUTH", instance.AuthMode)
	}

	return nil
}
