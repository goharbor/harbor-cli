package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	ScopeSelectors ImmutableSelector `json:"scope_selectors,omitempty"`
	TagSelectors ImmutableSelector `json:"tag_selectors"`
}

type ImmutableSelector struct {
	Decoration string `json:"decoration,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

func CreateImmutableView(createView *CreateView) {
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
				Description("Enter multiple comma separated repos,repo*,or **").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("pattern cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}