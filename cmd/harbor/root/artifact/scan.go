package artifact

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/scan"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ScanArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scan",
		Short:   "Scan an artifact",
		Long:    `Scan an artifact in Harbor Repository`,
		Example: `harbor artifact scan start <project>/<repository>/<reference>`,
	}

	cmd.AddCommand(
		StartScanArtifactCommand(),
		StopScanArtifactCommand(),
		// LogScanArtifactCommand(),
	)

	return cmd
}

func StartScanArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start a scan of an artifact",
		Long:    `Start a scan of an artifact in Harbor Repository`,
		Example: `harbor artifact scan start <project>/<repository>/<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				err = runStartScanArtifact(projectName, repoName, reference)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				err = runStartScanArtifact(projectName, repoName, reference)
			}

			if err != nil {
				log.Errorf("failed to start a scan of an artifact: %v", err)
			}
		},
	}
	return cmd
}

func runStartScanArtifact(projectName, repoName, reference string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	_, err := client.Scan.ScanArtifact(ctx, &scan.ScanArtifactParams{ProjectName: projectName, RepositoryName: repoName, Reference: reference})

	if err != nil {
		return err
	}

	log.Infof("Scan started successfully")

	return nil
}

func runStopScanArtifact(projectName, repoName, reference string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	_, err := client.Scan.StopScanArtifact(ctx, &scan.StopScanArtifactParams{ProjectName: projectName, RepositoryName: repoName, Reference: reference})

	if err != nil {
		return err
	}

	log.Infof("Scan stopped successfully")

	return nil
}

func StopScanArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stop a scan of an artifact",
		Long:    `Stop a scan of an artifact in Harbor Repository`,
		Example: `harbor artifact scan stop <project>/<repository>/<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				err = runStopScanArtifact(projectName, repoName, reference)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				err = runStopScanArtifact(projectName, repoName, reference)
			}

			if err != nil {
				log.Errorf("failed to stop a scan of an artifact: %v", err)
			}
		},
	}
	return cmd
}
