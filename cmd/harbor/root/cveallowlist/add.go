package cveallowlist

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/systemcve/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func AddCveAllowlistCommand() *cobra.Command {
	var opts update.UpdateView

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add cve allowlist",
		Long:  "Create allowlists of CVEs to ignore during vulnerability scanning",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			updateView := &update.UpdateView{
				CveId:      opts.CveId,
				IsExpire:   opts.IsExpire,
				ExpireDate: opts.ExpireDate,
			}

			err = updatecveView(updateView)
			if err != nil {
				log.Errorf("failed to add cveallowlist: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.IsExpire, "isexpire", "i", false, "Indicates whether the CVE entries should have an expiration date. Set to true to specify an expiration date")
	flags.StringVarP(&opts.CveId, "cveid", "n", "", "Comma-separated list of CVE IDs to be added to the allowlist")
	flags.StringVarP(&opts.ExpireDate, "expiredate", "d", "", "Specifies the expiration date for the CVE entries in the format 'YYYY-MM-DD'")

	return cmd
}

func updatecveView(updateView *update.UpdateView) error {
	if updateView == nil {
		updateView = &update.UpdateView{}
	}

	update.UpdateCveView(updateView)
	return api.UpdateSystemCve(*updateView)
}
