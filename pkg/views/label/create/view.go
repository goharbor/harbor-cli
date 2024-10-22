package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct{
	Name	string
	Color	string
	Description	string
	Scope	string
}

func CreateLabelView(createView *CreateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("name cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Color").
				Options(
					huh.NewOption("Black", "#000000"),
					huh.NewOption("Gray", "#737373"),
					huh.NewOption("White", "#FFFFFF"),
					huh.NewOption("Alto", "#DDDDDD"),
					huh.NewOption("Endeavour", "#0065AB"),
					huh.NewOption("Cerulean", "#0095D3"),
					huh.NewOption("Rose", "#F1428A"),
					huh.NewOption("Red", "#C92100"),
					huh.NewOption("Orange", "#F57600"),
					huh.NewOption("Yellow", "#FFDC0B"),
					huh.NewOption("Green", "#48960C"),
					huh.NewOption("Blue", "#343DAC"),
				).
				Value(&createView.Color).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("color cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}