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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/project/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListProjectCommand() *cobra.Command {
	var opts api.ListFlags
	var private bool
	var public bool
	var projects project.ListProjectsOK
	var err error
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list project",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if private && public {
				log.Fatal("Cannot specify both --private and --public flags")
			} else if private {
				opts.Public = false
				projects, err = api.ListProject(opts)
			} else if public {
				opts.Public = true
				projects, err = api.ListProject(opts)
			} else {
				projects, err = api.ListAllProjects(opts)
			}

			if err != nil {
				log.Fatalf("failed to get projects list: %v", err)
				return
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(projects, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListProjects(projects.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the project")
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.BoolVarP(&private, "private", "", false, "Show only private projects")
	flags.BoolVarP(&public, "public", "", false, "Show only public projects")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}
