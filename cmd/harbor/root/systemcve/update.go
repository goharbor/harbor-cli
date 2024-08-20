package systemcve

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/systemcve/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateSystemCveCommand() *cobra.Command {
	var opts update.UpdateView

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update systemcve allowlist",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			updateView := &update.UpdateView{
				CveId:      opts.CveId,
				IsExpire:   opts.IsExpire,
				ExpireDate: opts.ExpireDate,
			}

			err = updatecveView(updateView)
			if err != nil {
				log.Errorf("failed to update systemcve: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.IsExpire, "isexpire", "i", false, "Systemcve allowlist expire or not")
	flags.StringVarP(&opts.CveId, "cveid", "n", "", "CVE ids seperate with commas")
	flags.StringVarP(&opts.ExpireDate, "expiredate", "d", "", "If it expire,give Expiry date in the format MM/DD/YYYY")

	return cmd
}

func updatecveView(updateView *update.UpdateView) error {
	if updateView == nil {
		updateView = &update.UpdateView{}
	}

	update.UpdateCveView(updateView)
	return api.UpdateSystemCve(*updateView)
}
