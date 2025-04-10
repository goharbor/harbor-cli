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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/instance/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateInstanceCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create instance",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
				Name:        opts.Name,
				Vendor:      opts.Vendor,
				Description: opts.Description,
				Endpoint:    opts.Endpoint,
				Insecure:    opts.Insecure,
				Enabled:     opts.Enabled,
				AuthMode:    opts.AuthMode,
				AuthInfo:    opts.AuthInfo,
			}

			if opts.Name != "" && opts.Vendor != "" && opts.Endpoint != "" {
				err = api.CreateInstance(opts)
			} else {
				err = createInstanceView(createView)
			}

			if err != nil {
				log.Errorf("failed to create instance: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the instance")
	flags.StringVarP(&opts.Vendor, "provider", "p", "", "Provider for the instance")
	flags.StringVarP(&opts.Endpoint, "url", "u", "", "URL for the instance")
	flags.StringVarP(&opts.Description, "description", "", "", "Description of the instance")
	flags.BoolVarP(&opts.Insecure, "insecure", "i", true, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.BoolVarP(&opts.Enabled, "enable", "", true, "Whether it is enable or not")
	flags.StringVarP(&opts.AuthMode, "authmode", "a", "NONE", "Choosing different types of authentication method")

	return cmd
}

func createInstanceView(createView *create.CreateView) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateInstanceView(createView)
	return api.CreateInstance(*createView)
}
