package scanner

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/scanner/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update scanner",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var registrationID string

			if len(args) > 0 {
				registrationID = args[0]
			} else {
				registrationID = prompt.GetScannerIdFromUser()
			}

			if opts.Name == "" || opts.Description == "" || opts.Auth == "" || opts.AccessCredential == "" || opts.URL == "" {
				create.CreateScannerView(&opts)
			}

			err = api.UpdateScanner(registrationID, opts)

			if err != nil {
				log.Errorf("failed to update scanner: %v", err)
			} else {
				log.Infof("scanner updated successfully")
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the scanner")
	flags.StringVarP(&opts.Description, "des", "", "", "Description of the scanner")
	flags.StringVarP(&opts.Auth, "auth", "", "", "Authentication approach of the scanner [None|Basic|Bearer|X-ScannerAdapter-API-Key]")
	flags.StringVarP(&opts.AccessCredential, "cred", "", "", "HTTP Authorization header value sent with each request to the Scanner Adapter API")
	flags.StringVarP(&opts.URL, "url", "", "", "Base URL of the scanner adapter")
	flags.BoolVarP(&opts.Disabled, "disable", "", false, "Indicate whether the registration is enabled or not")
	flags.BoolVarP(&opts.SkipCertVerify, "skip", "", false, "Indicate if skip the certificate verification when sending HTTP requests")
	flags.BoolVarP(&opts.UseInternalAddr, "internal", "", false, "Indicate whether use internal registry addr for the scanner to pull content or not")

	return cmd
}
