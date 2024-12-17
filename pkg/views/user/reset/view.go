package reset

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type ResetView struct {
	OldPassword string
	NewPassword string
	Comfirm     string
}

func ResetUserView(resetView *ResetView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Old Password").
				EchoMode(huh.EchoModePassword).
				Value(&resetView.OldPassword).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("password cannot be empty or only spaces")
					}
					if err := utils.ValidatePassword(str); err != nil {
						return err
					}
					return nil
				}),
			huh.NewInput().
				Title("New Password").
				EchoMode(huh.EchoModePassword).
				Value(&resetView.NewPassword).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("password cannot be empty or only spaces")
					}
					if err := utils.ValidatePassword(str); err != nil {
						return err
					}
					return nil
				}),
			huh.NewInput().
				Title("Comfirm Password").
				EchoMode(huh.EchoModePassword).
				Value(&resetView.Comfirm).
				Validate(func(str string) error {
					if resetView.Comfirm != resetView.NewPassword {
						return errors.New("passwords do not match, please try again")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
