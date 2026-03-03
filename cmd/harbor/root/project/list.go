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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/project/list"
	log "github.com/sirupsen/logrus"
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
			log.Debug("Starting project list command")

			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}

			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			if private && public {
				return fmt.Errorf("Cannot specify both --private and --public flags")
			}

			var listFunc func(...api.ListFlags) (project.ListProjectsOK, error)
			if private {
				log.Debug("Using private project list function")
				opts.Public = false
				listFunc = api.ListProject
			} else if public {
				log.Debug("Using public project list function")
				opts.Public = true
				listFunc = api.ListProject
			} else {
				log.Debug("Using list all projects function")
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

			log.Debug("Fetching projects...")
			allProjects, err = fetchProjects(listFunc, opts)
			if err != nil {
				return fmt.Errorf("failed to get projects list: %v", utils.ParseHarborErrorMsg(err))
			}

			log.WithField("count", len(allProjects)).Debug("Number of projects fetched")
			if len(allProjects) == 0 {
				log.Info("No projects found")
				return nil
			}
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				log.WithField("output_format", formatFlag).Debug("Output format selected")
				err = utils.PrintFormat(allProjects, formatFlag)
				if err != nil {
					return err
				}
			} else {
				log.Debug("Listing projects using default view")
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
		log.Debug("Page size is 0, will fetch all pages")
		opts.PageSize = 100
		opts.Page = 1

		for {
			log.WithFields(log.Fields{
				"page":      opts.Page,
				"page_size": opts.PageSize,
			}).Debug("Fetching next page of projects")

			projects, err := listFunc(opts)
			if err != nil {
				return nil, err
			}

			log.WithField("fetched_count", len(projects.Payload)).Debug("Fetched projects from current page")
			allProjects = append(allProjects, projects.Payload...)

			if len(projects.Payload) < int(opts.PageSize) {
				log.Debug("Last page reached, stopping pagination")
				break
			}

			opts.Page++
		}
	} else {
		log.WithFields(log.Fields{
			"page":      opts.Page,
			"page_size": opts.PageSize,
		}).Debug("Fetching projects with user-defined pagination")

		projects, err := listFunc(opts)
		if err != nil {
			return nil, err
		}
		allProjects = projects.Payload
	}

	return allProjects, nil
}
