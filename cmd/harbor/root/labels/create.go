package labels

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/label/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateLabelCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use: "create",
		Short: "create label",
		Long: "crate label in harbor",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
				Name:	opts.Name,
				Color:	opts.Color,
				Scope:	opts.Scope,
				Description: opts.Description,
			}
			if opts.Name != "" && opts.Scope != "" {
				err = api.CreateLabels(opts)
			} else {
				err = createLabelView(createView)
			}

			if err != nil {
				log.Errorf("failed to create label: %v", err)
			}
			
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the label")
	flags.StringVarP(&opts.Color, "color", "c", "", "Color of the label.color is in hexadecimal value")
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "Scope of the label. eg- g(global), p(specific project)")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the label")
	
	return cmd
}

func createLabelView(createView *create.CreateView) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateLabelView(createView)
	return api.CreateLabels(*createView)
}

