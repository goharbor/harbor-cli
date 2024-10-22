package systeminfo

import (
	"context"
	"fmt"
	"os"
 

	"github.com/spf13/cobra"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
)

func GetCertCommand() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "cert",
		Short: "Get the default root certificate",
		Long:  `Download the default root certificate from the Harbor system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := utils.GetClient()
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			
			if outputFile == "" {
				outputFile = "harbor_root.crt"
			}
			file, err := os.Create(outputFile)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer file.Close()
			
			params := systeminfo.NewGetCertParams()
			_, err = client.Systeminfo.GetCert(context.Background(), params, file)
			if err != nil {
				switch err := err.(type) {
				case *systeminfo.GetCertNotFound:
					return fmt.Errorf("certificate not found. This could mean:\n" +
						"1. The certificate endpoint is not available in your Harbor version\n" +
						"2. Harbor is not configured to serve the root certificate\n" +
						"3. You may not have the necessary permissions\n" +
						"Please check your Harbor configuration and try again")
				default:
					return fmt.Errorf("failed to get certificate: %w", err)
				}
			}
			
			fmt.Printf("Certificate downloaded to %s\n", outputFile)
			return nil
		},
	}
	cmd.Flags().StringVar(&outputFile, "output", "", "Output file for the certificate (default: harbor_root.crt)")

	return cmd
}