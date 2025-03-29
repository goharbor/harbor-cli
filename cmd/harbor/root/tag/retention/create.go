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
		Use:   "create",
		Short: "Create a tag retention rule in a project",
		Long: `Create a tag retention rule for a project in Harbor to manage the lifecycle of image tags.

Tag retention rules help users automatically retain or delete specific tags based on 
defined criteria, reducing storage usage and improving repository maintenance.

⚠️ A user can create **up to 15 tag retention rules per project**.`,
		Example: `  # Retain tags matching 'release-*' at the project level
  harbor tag retention create --level project --action retain --taglist release-*

  # Delete untagged images at the repository level
  harbor retention create --level repository --action delete --tagdecoration untagged`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
				ScopeSelectors: create.RetentionSelector{
					Decoration: opts.ScopeSelectors.Decoration,
					Pattern:    opts.ScopeSelectors.Pattern,
				},
				TagSelectors: create.RetentionSelector{
					Decoration: opts.TagSelectors.Decoration,
					Pattern:    opts.TagSelectors.Pattern,
					Extras:     opts.TagSelectors.Extras,
				},
				Scope: create.RetentionPolicyScope{
					Level: opts.Scope.Level,
					Ref:   opts.Scope.Ref,
				},
				Template:  opts.Template,
				Params:    opts.Params,
				Action:    opts.Action,
				Algorithm: opts.Algorithm,
			}

			projectId := int32(prompt.GetProjectIDFromUser())
			err = createRetentionView(createView, projectId)

			if err != nil {
				log.Errorf("Failed to create retention tag rule: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ScopeSelectors.Decoration, "repodecoration", "", "", "Apply or exclude repositories from the rule")
	flags.StringVarP(&opts.ScopeSelectors.Pattern, "repolist", "", "", "Comma-separated list of repositories to apply/exclude")
	flags.StringVarP(&opts.TagSelectors.Decoration, "tagdecoration", "", "", "Apply or exclude specific tags from the rule")
	flags.StringVarP(&opts.TagSelectors.Pattern, "taglist", "", "", "Comma-separated list of tags to apply/exclude")
	flags.StringVarP(&opts.Scope.Level, "level", "", "project", "Scope of the retention policy: 'project' or 'repository'")
	flags.StringVarP(&opts.Action, "action", "", "retain", "Action to perform: 'retain' or 'delete'")
	flags.StringVarP(&opts.Algorithm, "algorithm", "", "or", "Rule combination method: 'or' or 'and'")

	return cmd
}

func createRetentionView(createView *create.CreateView, projectId int32) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateRetentionView(createView)
	return api.CreateRetention(*createView, projectId)
}
