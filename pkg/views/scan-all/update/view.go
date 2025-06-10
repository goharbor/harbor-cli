package update

import (
	"errors"
	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

func UpdateSchedule(cron *string) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the cron").
				Value(cron).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("cron cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
