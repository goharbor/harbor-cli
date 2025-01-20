package securityhub

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/securityhub/summary"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SummaryVulnerabilityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "The vulnerability summary of the system",
	}

	cmd.AddCommand(
		totalVulnerabilities(),
		MostDangerousArtifacts(),
		MostDangerousCVE(),
	)

	return cmd
}

func totalVulnerabilities() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total",
		Short: "Total Vulnerabilities",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			sum, err := api.GetTotalVulnerabilities(false, false)

			if err != nil {
				log.Fatalf("failed to get vulnerability summary: %v", err)
				return
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(sum, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				summary.GetTotalVulnerability(sum.Payload)
			}
		},
	}

	return cmd
}

func MostDangerousArtifacts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Top 5 Most Dangerous Artifacts",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			sum, err := api.GetTotalVulnerabilities(true, false)

			if err != nil {
				log.Fatalf("failed to get most dangerous artifacts: %v", err)
				return
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(sum, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				summary.ShowMostDangerousArtifacts(sum.Payload)
			}
		},
	}

	return cmd
}

func MostDangerousCVE() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cvelist",
		Short: "Top 5 Most Dangerous CVEs",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			sum, err := api.GetTotalVulnerabilities(false, true)

			if err != nil {
				log.Fatalf("failed to get most dangerous cvelist: %v", err)
				return
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(sum, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				summary.ShowMostDangerousCVE(sum.Payload)
			}
		},
	}

	return cmd
}
