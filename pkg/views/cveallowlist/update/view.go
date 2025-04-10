package update

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type UpdateView struct {
	CveId      string
	IsExpire   bool
	ExpireDate string
}

func UpdateCveView(updateView *UpdateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("CVE ID").
				Value(&updateView.CveId).
				Description("CVE IDs are separator by commas").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("cve id cannot be empty")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Expires").
				Value(&updateView.IsExpire).
				Affirmative("Date").
				Negative("never"),
		),
		huh.NewGroup(
			huh.NewInput().
				Validate(func(str string) error {
					if str == "" {
						return errors.New("ExpireDate cannot be empty")
					}
					return nil
				}).
				Description("Expire Date in the format YYYY/MM/DD").
				Title("Expire Date").
				Value(&updateView.ExpireDate),
		).WithHideFunc(func() bool {
			return !updateView.IsExpire
		}),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
