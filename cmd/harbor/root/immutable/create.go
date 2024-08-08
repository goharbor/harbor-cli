package immutable

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/immutable/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateImmutableCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use: "create",
		Short: "create immutable tag rule",
		Long: "create immutable tag rule to the project in harbor",
		Args: cobra.MaximumNArgs(1),
		Example: "harbor immutable create",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
				ScopeSelectors: create.ImmutableSelector{
					Decoration:	opts.ScopeSelectors.Decoration,
					Pattern:	opts.ScopeSelectors.Pattern,
				},
				TagSelectors: create.ImmutableSelector{
					Decoration:	opts.TagSelectors.Decoration,
					Pattern:	opts.TagSelectors.Pattern,
				},
			}
			if len(args) > 0 {
				err = createImmutableView(createView,args[0])
			} else {
				projectName := prompt.GetProjectNameFromUser()
				err = createImmutableView(createView,projectName)
			}

			if err != nil {
				log.Errorf("failed to create immutable tag rule: %v", err)
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ScopeSelectors.Decoration, "repo-decoration", "", "", "repository which either apply or exclude from the rule")
	flags.StringVarP(&opts.ScopeSelectors.Pattern, "repo-list", "", "", "list of repository to which to either apply or exclude from the rule")
	flags.StringVarP(&opts.TagSelectors.Decoration, "tag-decoration", "", "", "tags which either apply or exclude from the rule")
	flags.StringVarP(&opts.TagSelectors.Pattern, "tag-list", "", "", "list of tags to which to either apply or exclude from the rule")

	return cmd
}

func createImmutableView(createView *create.CreateView,projectName string) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateImmutableView(createView)
	return api.CreateImmutable(*createView,projectName)
}