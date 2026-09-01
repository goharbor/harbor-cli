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
	"log/slog"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/project/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListProjectCommand() *cobra.Command {
	var (
		opts        api.ListFlags
		private     bool
		public      bool
		allProjects []*models.Project
		err         error
		// For querying, opts.Q
		fuzzy  []string
		match  []string
		ranges []string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Debug("Starting project list command")
			if err := utils.ValidatePagination(opts.Page, opts.PageSize); err != nil {
				return err
			}

			if private && public {
				return fmt.Errorf("Cannot specify both --private and --public flags")
			}

			var listFunc func(...api.ListFlags) (project.ListProjectsOK, error)
			if private {
				slog.Debug("Using private project list function")
				opts.Public = false
				listFunc = api.ListProject
			} else if public {
				slog.Debug("Using public project list function")
				opts.Public = true
				listFunc = api.ListProject
			} else {
				slog.Debug("Using list all projects function")
				listFunc = api.ListAllProjects
			}

			if len(fuzzy) != 0 || len(match) != 0 || len(ranges) != 0 { // Only Building Query if a param exists
				q, qErr := utils.BuildQueryParam(fuzzy, match, ranges,
					[]string{"name", "project_id", "public", "creation_time", "owner_id"},
				)
				if qErr != nil {
					return qErr
				}

				opts.Q = q
			}

			slog.Debug("Fetching projects...")
			allProjects, err = fetchProjects(listFunc, opts)
			if err != nil {
				return fmt.Errorf("failed to get projects list: %v", utils.ParseHarborErrorMsg(err))
			}

			slog.Debug("Number of projects fetched", "count", len(allProjects))
			if len(allProjects) == 0 {
				fmt.Println("No projects found")
				return nil
			}
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				slog.Debug("Output format selected", "output_format", formatFlag)
				err = utils.PrintFormat(allProjects, formatFlag)
				if err != nil {
					return err
				}
			} else {
				slog.Debug("Listing projects using default view")
				list.ListProjects(allProjects)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the project")
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 0, "Size of per page (0 to fetch all)")
	flags.BoolVarP(&private, "private", "", false, "Show only private projects")
	flags.BoolVarP(&public, "public", "", false, "Show only public projects")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")
	flags.StringSliceVar(&fuzzy, "fuzzy", nil, "Fuzzy match filter (key=value)")
	flags.StringSliceVar(&match, "match", nil, "exact match filter (key=value)")
	flags.StringSliceVar(&ranges, "range", nil, "range filter (key=min~max)")

	return cmd
}

func fetchProjects(listFunc func(...api.ListFlags) (project.ListProjectsOK, error), opts api.ListFlags) ([]*models.Project, error) {
	var allProjects []*models.Project
	if opts.PageSize == 0 {
		slog.Debug("Page size is 0, will fetch all pages")
		opts.PageSize = 100
		opts.Page = 1

		for {
			slog.Debug("Fetching next page of projects", "page", opts.Page, "page_size", opts.PageSize)

			projects, err := listFunc(opts)
			if err != nil {
				return nil, err
			}

			slog.Debug("Fetched projects from current page", "fetched_count", len(projects.Payload))
			allProjects = append(allProjects, projects.Payload...)

			if len(projects.Payload) < int(opts.PageSize) {
				slog.Debug("Last page reached, stopping pagination")
				break
			}

			opts.Page++
		}
	} else {
		slog.Debug("Fetching projects with user-defined pagination", "page", opts.Page, "page_size", opts.PageSize)

		projects, err := listFunc(opts)
		if err != nil {
			return nil, err
		}
		allProjects = projects.Payload
	}

	return allProjects, nil
}
