package usergroup

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/api"
	create "github.com/goharbor/harbor-cli/pkg/views/usergroup/create"
	"github.com/spf13/cobra"
)

type ErrorResponse struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func UserGroupCreateCommand() *cobra.Command {
	var opts create.CreateUserGroupInput

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create user group",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Getting Vals
			err := create.CreateUserGroupView(&opts)
			if err != nil {
				return err
			}

			fmt.Printf("Creating user group with name: %s, type: %d%s\n", opts.GroupName, opts.GroupType, opts.LDAPGroupDN)
			err = api.CreateUserGroup(opts.GroupName, opts.GroupType, opts.LDAPGroupDN)
			if err != nil {
				return formatError(err)
			}

			fmt.Printf("User group '%s' created successfully\n", opts.GroupName)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.GroupName, "name", "n", "", "Group name")
	flags.Int64VarP(&opts.GroupType, "type", "t", 0, "Group type")
	flags.StringVarP(&opts.LDAPGroupDN, "ldap-dn", "l", "", "The DN of the LDAP group if group type is 1 (LDAP group)")

	return cmd
}

func formatError(err error) error {
	errStr := err.Error()
	if strings.Contains(errStr, "conflict:") {
		var errResp ErrorResponse
		jsonStr := strings.TrimPrefix(errStr, "conflict: ")
		if err := json.Unmarshal([]byte(jsonStr), &errResp); err == nil {
			if len(errResp.Errors) > 0 {
				return fmt.Errorf("%s", errResp.Errors[0].Message)
			}
		}
	}
	return fmt.Errorf("failed to create user group: %v", err)
}
