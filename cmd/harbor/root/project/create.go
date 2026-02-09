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
	"io"
	"os"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CreateProjectCommand creates a new `harbor create project` command
type ProjectCreator interface {
	CreateProject(opts create.CreateView) error
	FillProjectView(createView *create.CreateView) error
}
type DefaultProjectCreator struct{}

func (d *DefaultProjectCreator) CreateProject(opts create.CreateView) error {
	return api.CreateProject(opts)
}
func (d *DefaultProjectCreator) FillProjectView(createView *create.CreateView) error {
	return create.CreateProjectView(createView)
}

func CreateProject(w io.Writer, projectCreator ProjectCreator, opts *create.CreateView, args []string) error {
	var ProjectName string
	var createView *create.CreateView
	if len(args) > 0 {
		opts.ProjectName = args[0]
	}

	if opts.ProxyCache && opts.RegistryID == "" {
		return fmt.Errorf("proxy cache selected but no registry ID provided. Use --registry-id")
	}

	if !opts.ProxyCache && opts.RegistryID != "" {
		return fmt.Errorf("registry ID should only be provided when proxy-cache is enabled")
	}

	if opts.ProjectName == "" || opts.StorageLimit == "" {
		log.Debug("Switching to interactive view...")
		createView = &create.CreateView{
			ProjectName:  opts.ProjectName,
			Public:       opts.Public,
			RegistryID:   opts.RegistryID,
			StorageLimit: opts.StorageLimit,
			ProxyCache:   opts.ProxyCache,
		}

		err := fillCreateView(projectCreator, createView)
		if err != nil {
			return fmt.Errorf("Failed to get the required params to create project:%w", err)
		}
	} else {
		createView = opts
	}

	ProjectName = createView.ProjectName
	if err := projectCreator.CreateProject(*createView); err != nil {
		return fmt.Errorf("failed to create project: %v", utils.ParseHarborErrorMsg(err))
	}
	fmt.Fprintf(w, "Project '%s' created successfully\n", ProjectName)
	return nil
}
func CreateProjectCommand() *cobra.Command {
	opts := &create.CreateView{}

	cmd := &cobra.Command{
		Use:   "create [project name]",
		Short: "create project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return CreateProject(os.Stdout, &DefaultProjectCreator{}, opts, args)
		}}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.Public, "public", "", false, "Project is public or private")
	flags.StringVarP(&opts.RegistryID, "registry-id", "", "", "ID of referenced registry when creating the proxy cache project")
	flags.StringVarP(&opts.StorageLimit, "storage-limit", "", "", "Storage quota of the project")
	flags.BoolVarP(&opts.ProxyCache, "proxy-cache", "", false, "Whether the project is a proxy cache project")

	return cmd
}

func fillCreateView(projectCreator ProjectCreator, createView *create.CreateView) error {
	if createView == nil {
		createView = &create.CreateView{
			ProjectName:  "",
			Public:       false,
			RegistryID:   "",
			StorageLimit: "-1",
		}
	}
	err := projectCreator.FillProjectView(createView)
	return err
}
