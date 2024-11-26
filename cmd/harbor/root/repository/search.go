package repository

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/search"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "search repository based on their names",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			repo, err := api.SearchRepository(args[0])
			if err != nil {
				log.Fatalf("failed to get repositories: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(repo)
				return
			}

			search.SearchRepositories(repo.Payload.Repository)
		},
	}
	return cmd
}
