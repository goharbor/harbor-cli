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
	"strings"

	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"

	"github.com/goharbor/harbor-cli/pkg/views/user/create"
	"github.com/spf13/cobra"
)

type UserCreator interface {
	FillUser(opts *create.CreateView)
	UserCreate(opts create.CreateView) error
}
type DefaultUserCreator struct{}

func (d *DefaultUserCreator) FillUser(opts *create.CreateView) {
	create.CreateUserView(opts)
}

func (d *DefaultUserCreator) UserCreate(opts create.CreateView) error {
	return api.CreateUser(opts)
}
func CreateUser(opts *create.CreateView, userCreator UserCreator) {
	var err error

	if opts.Email == "" || opts.Realname == "" || opts.Password == "" || opts.Username == "" {
		userCreator.FillUser(opts)
	}

	err = userCreator.UserCreate(*opts)

	if err != nil {
		if isUnauthorizedError(err) {
			log.Error("Permission denied: Admin privileges are required to execute this command.")
		} else {
			log.Errorf("failed to create user: %v", err)
		}
	}
}
func UserCreateCmd() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create user",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			d := &DefaultUserCreator{}
			CreateUser(&opts, d)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Email, "email", "", "", "Email")
	flags.StringVarP(&opts.Realname, "realname", "", "", "Realname")
	flags.StringVarP(&opts.Comment, "comment", "", "", "Comment")
	flags.StringVarP(&opts.Password, "password", "", "", "Password")
	flags.StringVarP(&opts.Username, "username", "", "", "Username")

	return cmd
}
func isUnauthorizedError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "403")
}
