package views

import (
	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

func ConfirmElevation() (bool, error) {
	var confirm bool

	err := huh.NewConfirm().
		Title("Are you sure to elevate the user to admin role?").
		Affirmative("Yes").
		Negative("No").
		Value(&confirm).Run()
	if err != nil {
		log.Fatal(err)
	}

	return confirm, nil
}
