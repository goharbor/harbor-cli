package repository

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// NewUpdateRepositoryCommand creates a new `harbor update repository` command
func UpdateRepositoryCommand() *cobra.Command {
	var opts config.UpdateRepositoryOptions

	cmd := &cobra.Command{
		Use:   "repository [ID & PROJECT_ID]",
		Short: "update repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunUpdateRepository(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ProjectName, "projectname", "", "", "Name of the project")
	flags.Int64VarP(&opts.Repository.ProjectId, "projectid", "", 0, "Id of the project")
	flags.StringVarP(&opts.RepositoryName, "reponame", "", "", "Name of the old repository")
	flags.Int64VarP(&opts.Repository.Id, "repoid", "", 0, "Id of the old repository")
	flags.StringVarP(&opts.Repository.Description, "description", "", "", "Description of the updated repository")
	flags.StringVarP(&opts.Repository.Name, "newreponame", "", "", "Name of the updated repository")
	flags.Int64VarP(&opts.Repository.ArtifactCount, "artifactcount", "", 0, "No of artifact present in the old repository")
	flags.Int64VarP(&opts.Repository.PullCount, "pullcount", "", 0, "No of times artifact pulled from the old repository")
	cmd.MarkFlagRequired("projname")
	cmd.MarkFlagRequired("projectname")
	cmd.MarkFlagRequired("newreponame")
	cmd.MarkFlagRequired("repoid")
	cmd.MarkFlagRequired("projectid")

	return cmd
}
