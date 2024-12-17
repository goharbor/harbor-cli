package artifact

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ScanArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scan",
		Short:   "Scan an artifact",
		Long:    `Scan an artifact in Harbor Repository`,
		Example: `harbor artifact scan start <project>/<repository>@<reference>`,
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
		Example: `harbor artifact scan start <project>/<repository>@<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				err = api.StartScanArtifact(projectName, repoName, reference)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				reference := prompt.GetReferenceFromUser(repoName, projectName)
				err = api.StartScanArtifact(projectName, repoName, reference)
			}
			if err != nil {
				log.Errorf("failed to start scan of artifact: %v", err)
			}
		},
	}
	return cmd
}

func StopScanArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stop a scan of an artifact",
		Long:    `Stop a scan of an artifact in Harbor Repository`,
		Example: `harbor artifact scan stop <project>/<repository>@<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				err = api.StopScanArtifact(projectName, repoName, reference)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				reference := prompt.GetReferenceFromUser(repoName, projectName)
				err = api.StopScanArtifact(projectName, repoName, reference)
			}
			if err != nil {
				log.Errorf("failed to stop scan of artifact: %v", err)
			}
		},
	}
	return cmd
}
