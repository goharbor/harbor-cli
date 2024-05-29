package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

func CreateTagView(tagName *string) {
	theme := huh.ThemeCharm()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Tag Name").
				Value(tagName).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("project name cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
