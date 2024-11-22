package artifact

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	artifactViews "github.com/goharbor/harbor-cli/pkg/views/artifact/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListArtifactCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list artifacts within a repository",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var resp artifact.ListArtifactsOK

			if len(args) > 0 {
				projectName, repoName := utils.ParseProjectRepo(args[0])
				resp, err = api.ListArtifact(projectName, repoName, opts)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				resp, err = api.ListArtifact(projectName, repoName, opts)
			}

			if err != nil {
				log.Errorf("failed to list artifacts: %v", err)
			}
			artifactViews.ListArtifacts(resp.Payload)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "p", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "n", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "s", "", "Sort the resource list in ascending or descending order")

	return cmd
}
