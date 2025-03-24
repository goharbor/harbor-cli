package retention

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/retention/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateRetentionCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use: 	"create",
		Short: 	"create retention tag rule",
		Long: 	"create retention tag rule to the project in harbor",
		Example: "harbor retention create",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
				ScopeSelectors: create.RetentionSelector{
					Decoration:	opts.ScopeSelectors.Decoration,
					Pattern:	opts.ScopeSelectors.Pattern,
				},
				TagSelectors: create.RetentionSelector{
					Decoration:	opts.TagSelectors.Decoration,
					Pattern:	opts.TagSelectors.Pattern,
					Extras:		opts.TagSelectors.Extras,
				},
				Scope: create.RetentionPolicyScope{
					Level: opts.Scope.Level,
					Ref: opts.Scope.Ref,
				},
				Template: opts.Template,
				Params: opts.Params,
				Action: opts.Action,
				Algorithm: opts.Algorithm,
			}

			projectId := int32(prompt.GetProjectIDFromUser())
			err = createRetentionView(createView,projectId)

			if err != nil {
				log.Errorf("failed to create retention tag rule: %v", err)
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ScopeSelectors.Decoration, "repodecoration", "", "", "repository which either apply or exclude from the rule")
	flags.StringVarP(&opts.ScopeSelectors.Pattern, "repolist", "", "", "list of repository to which to either apply or exclude from the rule")
	flags.StringVarP(&opts.TagSelectors.Decoration, "tagdecoration", "", "", "tags which either apply or exclude from the rule")
	flags.StringVarP(&opts.TagSelectors.Pattern, "taglist", "", "", "list of tags to which to either apply or exclude from the rule")
	flags.StringVarP(&opts.Scope.Level,"level","","project","scope of retention policy")
	flags.StringVarP(&opts.Action,"action","","retain","Action of the retention policy")
	flags.StringVarP(&opts.Algorithm,"algorithm","","or","Algorithm of retention policy")

	return cmd
}

func createRetentionView(createView *create.CreateView,projectId int32) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateRetentionView(createView)
	return api.CreateRetention(*createView,projectId)
}