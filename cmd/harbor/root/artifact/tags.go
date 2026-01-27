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
package artifact

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/artifact/tags/create"
	"github.com/goharbor/harbor-cli/pkg/views/artifact/tags/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ArtifactTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tags",
		Short:   "Manage tags of an artifact",
		Example: ` harbor artifact tags list <project>/<repository>/<reference>`,
	}

	cmd.AddCommand(
		ListTagsCmd(),
		DeleteTagsCmd(),
		CreateTagsCmd(),
	)

	return cmd
}

func CreateTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a tag of an artifact",
		Example: `harbor artifact tags create <project>/<repository>/<reference> <tag>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, repoName, reference string
			var tagName string
			if len(args) > 0 {
				projectName, repoName, reference, err = utils.ParseProjectRepoReference(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo/reference: %w", err)
				}
				if len(args) < 2 {
					return fmt.Errorf("missing tag argument")
				}
				tagName = args[1]
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %w", err)
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
				reference = prompt.GetReferenceFromUser(repoName, projectName)
				create.CreateTagView(&tagName)
			}
			err = api.CreateTag(projectName, repoName, reference, tagName)
			if err != nil {
				return fmt.Errorf("failed to create tag: %w", err)
			}
			return nil
		},
	}

	return cmd
}

func ListTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tags of an artifact",
		Example: `harbor artifact tags list <project>/<repository>/<reference>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var tags *artifact.ListTagsOK
			var projectName, repoName, reference string

			if len(args) > 0 {
				projectName, repoName, reference, err = utils.ParseProjectRepoReference(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo/reference: %w", err)
				}
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %w", err)
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
				if repoName == "" {
					return fmt.Errorf("repository name is required")
				}
				reference = prompt.GetReferenceFromUser(repoName, projectName)
			}

			tags, err = api.ListTags(projectName, repoName, reference)

			if err != nil {
				return fmt.Errorf("failed to list tags: %w", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(tags, FormatFlag)
				if err != nil {
					return fmt.Errorf("failed to print format: %w", err)
				}
			} else {
				list.ListTags(tags.Payload)
			}
			return nil
		},
	}

	return cmd
}

func DeleteTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a tag of an artifact",
		Example: `harbor artifact tags delete <project>/<repository>/<reference> <tag>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, repoName, reference string
			var tagName string
			if len(args) > 0 {
				projectName, repoName, reference, err = utils.ParseProjectRepoReference(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo/reference: %w", err)
				}
				if len(args) < 2 {
					return fmt.Errorf("missing tag argument")
				}
				tagName = args[1]
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %w", err)
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
				reference = prompt.GetReferenceFromUser(repoName, projectName)
				tagName = prompt.GetTagFromUser(repoName, projectName, reference)
			}
			err = api.DeleteTag(projectName, repoName, reference, tagName)
			if err != nil {
				return fmt.Errorf("failed to delete tag: %w", err)
			}
			return nil
		},
	}

	return cmd
}
