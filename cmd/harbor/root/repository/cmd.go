package repository

import "github.com/spf13/cobra"

func Repository() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage repositories",
		Long:  `Manage repositories in Harbor context`,
	}
	cmd.AddCommand(
		ListRepositoryCommand(),
		RepoViewCmd(),
		RepoDeleteCmd(),
		SearchRepoCmd(),
	)

	return cmd

}
