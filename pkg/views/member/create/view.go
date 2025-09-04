package create

import (
	"errors"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

type MemberUser struct {
	UserID   int
	Username string
}

type MemberGroup struct {
	ID          int
	GroupName   string
	GroupType   int
	LdapGroupDN string
}

type CreateView struct {
	ProjectName string
	RoleID      int
	RoleName    string
	MemberUser  *models.UserEntity
	MemberGroup *models.UserGroup
}

// map role names to role ids
var RoleOptions = map[string]int{
	"Admin":        1,
	"Developer":    2,
	"Guest":        3,
	"Maintainer":   4,
	"LimitedGuest": 5,
}

func CreateMemberView(createView *CreateView) {
	roleOptions := []string{"Project Admin", "Developer", "Guest", "Maintainer", "Limited Guest"}
	var roleSelectOptions []huh.Option[int]
	for id, name := range roleOptions {
		roleSelectOptions = append(roleSelectOptions, huh.NewOption(name, id))
	}

	groupOptions := []string{"None", "LDAP group", "HTTP group", "OIDC group"}
	var groupSelectOptions []huh.Option[int]
	for id, name := range groupOptions {
		groupSelectOptions = append(groupSelectOptions, huh.NewOption(name, id))
	}

	groups := []*huh.Group{}

	// Show role select only if not already given
	if createView.RoleID == 0 && createView.RoleName == "" {
		groups = append(groups,
			huh.NewGroup(
				huh.NewSelect[int]().
					Description("Select a Role").
					Title("Role").
					Value(&createView.RoleID).
					Options(roleSelectOptions...),
			))
	}

	if createView.MemberUser.UserID == 0 && createView.MemberUser.Username == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(&createView.MemberUser.Username).
				Validate(func(str string) error { // TODO: Add username checking
					if str == "" {
						return errors.New("Username and UserID cannot both be empty.")
					}

					return nil
				}),
		))
	}

	// always show optioal group name
	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title("Group Name (optional)").
			Value(&createView.MemberGroup.GroupName).
			Validate(func(str string) error {
				return nil
			}),
	))

	// if groupname is populated, show groupType
	if createView.MemberGroup.GroupName != "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Group Name (optional)").
				Value(&createView.MemberGroup.GroupName).
				Validate(func(str string) error {
					return nil
				}),
		), huh.NewGroup(
			huh.NewInput().
				Title("DN of LDAP group (optional)").
				Value(&createView.MemberGroup.GroupName).
				Validate(func(str string) error {
					return nil
				}),
		))
	}

	theme := huh.ThemeCharm()
	err := huh.NewForm(groups...).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}
