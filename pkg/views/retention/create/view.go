package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Action 			string
	ScopeSelectors 	RetentionSelector 
	TagSelectors 	RetentionSelector 
	Template 		string 
	Algorithm		string
	Params 			ParamsValue 
	Scope 			RetentionPolicyScope 
}

type RetentionSelector struct {
	Decoration 	string 
	Pattern 	string 
	Extras 		string 
	Kind 		string 
}

type ParamsValue struct {
	Name 	string 
	Value 	string 
}

type RetentionPolicyScope struct {
	Level 	string 
	Ref 	int64 
}

func CreateRetentionView(createView *CreateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("\nFor the repositories\n").
				Options(
					huh.NewOption("matching", "repoMatches"),
					huh.NewOption("excluding", "repoExcludes"),
				).Value(&createView.ScopeSelectors.Decoration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("decoration cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("List of repositories").
				Value(&createView.ScopeSelectors.Pattern).
				Description("Enter multiple comma separated repos,repo*,or **").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("pattern cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Tags\n").
				Options(
					huh.NewOption("matching", "matches"),
					huh.NewOption("excluding", "excludes"),
				).Value(&createView.TagSelectors.Decoration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("decoration cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("List of Tags").
				Value(&createView.TagSelectors.Pattern).
				Description("Enter multiple comma separated tags, tag*, or **.").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("pattern cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Untagged Artifacts\n").
				Description("Include or exclude all untagged artifacts by selecting true or false").
				Options(
					huh.NewOption("true", "{\"untagged\":true}"),
					huh.NewOption("false", "{\"untagged\":false}"),
				).Value(&createView.TagSelectors.Extras).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("this field cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("\nSelect the condition of retain\n").
				Options(
					huh.NewOption("retain the most recently pushed # artifacts", "latestPushedK"),
					huh.NewOption("retain the most recently pulled # artifacts", "latestPulledN"),
					huh.NewOption("retain the artifacts pushed within the last # days", "nDaysSinceLastPush"),
					huh.NewOption("retain the artifacts pulled within the last # days", "nDaysSinceLastPull"),
					huh.NewOption("retain always", "always"),
				).Value(&createView.Template).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("this field cannot be empty")
					}
					return nil
				}),
			),
		huh.NewGroup(
			huh.NewInput().
			Title("Count").
			Value(&createView.Params.Value).
			Description("Enter the number of artifact count").
			Validate(func(str string) error {
				if str == "" {
					return errors.New("count cannot be empty")
				}
				return nil
			}),
		).WithHideFunc(func() bool {
			return createView.Template == "always" || createView.Template == "nDaysSinceLastPush" || createView.Template == "nDaysSinceLastPull"
		}),
		huh.NewGroup(
			huh.NewInput().
			Title("Days").
			Value(&createView.Params.Value).
			Description("Enter the number of days").
			Validate(func(str string) error {
				if str == "" {
					return errors.New("days cannot be empty")
				}
				return nil
			}),
		).WithHideFunc(func() bool {
			return createView.Template == "always" || createView.Template == "latestPulledN" || createView.Template == "latestPushedK"
		}),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}