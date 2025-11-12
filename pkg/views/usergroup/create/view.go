package create

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

type CreateUserGroupInput struct {
	GroupName   string
	GroupType   int64
	LDAPGroupDN string
}

func CreateUserGroupView(opts *CreateUserGroupInput) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter new group name").
				Value(&opts.GroupName),
			huh.NewSelect[int64]().
				Title("Select group type").
				Options(
					huh.NewOption("LDAP", int64(1)),
					huh.NewOption("HTTP", int64(2)),
					huh.NewOption("OIDC", int64(3)),
				).
				Value(&opts.GroupType),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("form input error: %v", err)
	}

	// If LDAP
	if opts.GroupType == 1 {
		err := huh.NewInput().
			Title("LDAP Group DN").
			Value(&opts.LDAPGroupDN).
			Run()
		if err != nil {
			return fmt.Errorf("form input error: %v", err)
		}
	}

	return nil
}
