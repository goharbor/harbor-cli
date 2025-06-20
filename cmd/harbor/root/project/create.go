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
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CreateProjectCommand creates a new `harbor create project` command
func CreateProjectCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create [project name]",
		Short: "create project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var ProjectName string
			if len(args) > 0 {
				opts.ProjectName = args[0]
			}

			if opts.ProxyCache && opts.RegistryID == "" {
				return fmt.Errorf("proxy cache selected but no registry ID provided. Use --registry-id")
			}

			if !opts.ProxyCache && opts.RegistryID != "" {
				return fmt.Errorf("registry ID should only be provided when proxy-cache is enabled")
			}

			if opts.ProjectName != "" && opts.StorageLimit != "" {
				log.Debug("Attempting to create project using flags...")
				err = api.CreateProject(opts)
				ProjectName = opts.ProjectName
			} else {
				log.Debug("Switching to interactive view...")
				createView := &create.CreateView{
					ProjectName:  opts.ProjectName,
					Public:       opts.Public,
					RegistryID:   opts.RegistryID,
					StorageLimit: opts.StorageLimit,
					ProxyCache:   opts.ProxyCache,
				}

				err = createProjectView(createView)
				ProjectName = createView.ProjectName
			}

			if err != nil {
				return fmt.Errorf("failed to create project: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Printf("Project '%s' created successfully\n", ProjectName)
			return nil
		}}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.Public, "public", "", false, "Project is public or private")
	flags.StringVarP(&opts.RegistryID, "registry-id", "", "", "ID of referenced registry when creating the proxy cache project")
	flags.StringVarP(&opts.StorageLimit, "storage-limit", "", "", "Storage quota of the project")
	flags.BoolVarP(&opts.ProxyCache, "proxy-cache", "", false, "Whether the project is a proxy cache project")

	return cmd
}

func createProjectView(createView *create.CreateView) error {
	if createView == nil {
		createView = &create.CreateView{
			ProjectName:  "",
			Public:       false,
			RegistryID:   "",
			StorageLimit: "-1",
		}
	}

	err := create.CreateProjectView(createView)
	if err != nil {
		return err
	}
	return api.CreateProject(*createView)
}
