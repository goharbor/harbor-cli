package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	ProjectName  string
	Public       bool
	RegistryID   string
	StorageLimit string
}

func CreateProjectView(createView *CreateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Value(&createView.ProjectName).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("project name cannot be empty")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Public").
				Value(&createView.Public).
				Affirmative("yes").
				Negative("no"),
			huh.NewInput().
				Title("Registry ID").
				Value(&createView.RegistryID).
				Validate(func(str string) error {
					// Assuming RegistryID is an int64
					return nil
				}),
			huh.NewInput().
				Title("Storage Limit").
				Value(&createView.StorageLimit).
				Validate(func(str string) error {
					// Assuming StorageLimit is an int64
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
