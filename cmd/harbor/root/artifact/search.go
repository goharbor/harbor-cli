package artifact

import (
	"fmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	artifactViews "github.com/goharbor/harbor-cli/pkg/views/artifact/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func SearchArtifacts() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search container artifacts (images, charts, etc.) in a Harbor repository",
		Long: `List artifacts (e.g., container images, charts) within a given Harbor repository. 
Search is based on matching tags and artifact types (e.g., container, images, charts)

Examples:
  harbor-cli artifact search project/repo:tag               
  harbor-cli artifact search project/repo:tag --type IMAGE
`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}
			searchTerm := args[0]
			project, repoAndTag, found := strings.Cut(searchTerm, "/")
			if !found {
				return fmt.Errorf("invalid search term: %s", searchTerm)
			}
			repository, searchTag, found := strings.Cut(repoAndTag, ":")
			if !found {
				return fmt.Errorf("invalid search term: %s", searchTerm)
			}

			artifacts, err := api.ListArtifact(project, repository, opts)
			if err != nil {
				return fmt.Errorf("failed to list artifacts: %v", utils.ParseHarborErrorMsg(err))
			}

			var matching []*models.Artifact

			for _, af := range artifacts.Payload {
				found := false
				for _, tag := range af.Tags {
					if tag.Name == searchTag {
						found = true
					}
				}
				if found {
					matching = append(matching, af)
				}
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(matching, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				artifactViews.ListArtifacts(matching)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "p", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "n", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "s", "", "Sort the resource list in ascending or descending order")

	return cmd
}
