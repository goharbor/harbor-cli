package update

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	log "github.com/sirupsen/logrus"
)

type UpdateView struct {
	CreationTime strfmt.DateTime    `json:"creation_time,omitempty"`
	Description  string             `json:"description,omitempty"`
	Disable      bool               `json:"disable,omitempty"`
	Duration     int64              `json:"duration,omitempty"`
	Editable     bool               `json:"editable"`
	ExpiresAt    int64              `json:"expires_at,omitempty"`
	ID           int64              `json:"id,omitempty"`
	Level        string             `json:"level,omitempty"`
	Name         string             `json:"name,omitempty"`
	Permissions  []*RobotPermission `json:"permissions"`
	Secret       string             `json:"secret,omitempty"`
	UpdateTime   strfmt.DateTime    `json:"update_time,omitempty"`
}

type RobotPermission struct {
	Access    []*models.Access `json:"access"`
	Kind      string           `json:"kind,omitempty"`
	Namespace string           `json:"namespace,omitempty"`
}

type Access struct {
	Action   string `json:"action,omitempty"`
	Effect   string `json:"effect,omitempty"`
	Resource string `json:"resource,omitempty"`
}

func UpdateRobotView(updateView *UpdateView) {
	var duration string
	duration = strconv.FormatInt(updateView.Duration, 10)

	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Description").
				Value(&updateView.Description),
			huh.NewInput().
				Title("Expiration").
				Value(&duration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Expiration cannot be empty")
					}
					dur, err := strconv.ParseInt(str, 10, 64)
					if err != nil {
						return errors.New("invalid expiration time: Enter expiration time in days")
					}
					updateView.Duration = dur
					return nil
				}),
			huh.NewConfirm().
				Title("Disable").
				Value(&updateView.Disable).
				Affirmative("yes").
				Negative("no"),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}
