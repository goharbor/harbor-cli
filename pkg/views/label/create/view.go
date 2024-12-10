package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Name        string
	Color       string
	Description string
	Scope       string
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
					huh.NewOption("White", "#FFFFFF"),
					huh.NewOption("Black", "#000000"),
					huh.NewOption("Jet Grey", "#61717D"),
					huh.NewOption("Grey", "#737373"),
					huh.NewOption("Spicy Pink", "#80746D"),
					huh.NewOption("Cadet Blue", "#A9B6BE"),
					huh.NewOption("Alto", "#DDDDDD"),
					huh.NewOption("Silk", "#BBB3A9"),
					huh.NewOption("Endeavour", "#0065AB"),
					huh.NewOption("Sapphire", "#343DAC"),
					huh.NewOption("Violet", "#781DA0"),
					huh.NewOption("Jazzberry Jam", "#9B0D54"),
					huh.NewOption("Blue", "#0095D3"),
					huh.NewOption("Purple", "#9DA3DB"),
					huh.NewOption("Bright Lavender", "#BE90D6"),
					huh.NewOption("Rose", "#F1428A"),
					huh.NewOption("Navy Green", "#1D5100"),
					huh.NewOption("Dark Aqua", "#006668"),
					huh.NewOption("Peacock Blue", "#006690"),
					huh.NewOption("Regal Blue", "#004A70"),
					huh.NewOption("Green", "#48960C"),
					huh.NewOption("Cyan", "#00AB9A"),
					huh.NewOption("Cerulean", "#00B7D6"),
					huh.NewOption("Nice Blue", "#0081A7"),
					huh.NewOption("Red", "#C92100"),
					huh.NewOption("Thunderbird", "#CD3517"),
					huh.NewOption("Rust Orange", "#C25400"),
					huh.NewOption("Yellow Brown", "#D28F00"),
					huh.NewOption("Radical Red", "#F52F52"),
					huh.NewOption("Reddish Orange", "#FF5501"),
					huh.NewOption("Orange", "#F57600"),
					huh.NewOption("Yellow", "#FFDC0B"),
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
