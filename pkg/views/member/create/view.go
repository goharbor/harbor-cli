package create

import (
	"errors"
	"log"
	"strconv"

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
	ProjectNameOrID string
	RoleID          int
	MemberUser      *models.UserEntity
	MemberGroup     *models.UserGroup
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

	var (
		groupType int
		userID    string
		groupID   string
	)

	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Description("Select a Role").
				Title("Role").
				Value(&createView.RoleID).
				Options(roleSelectOptions...),

			huh.NewInput().
				Title("User ID").
				Value(&userID).
				Validate(func(str string) error {
					createView.MemberUser.UserID, _ = strconv.ParseInt(str, 10, 64)
					return nil
				}),
			huh.NewInput().
				Title("Username").
				Value(&createView.MemberUser.Username).
				Validate(func(str string) error {
					if userID == "" && str == "" {
						return errors.New("Username and UserID cannot both be empty.")
					}
					return nil
				}),
			huh.NewInput().
				Title("Group ID").
				Value(&groupID).
				Validate(func(str string) error {
					createView.MemberGroup.ID, _ = strconv.ParseInt(str, 10, 64)
					return nil
				}),
			huh.NewInput().
				Title("Group Name").
				Value(&createView.MemberGroup.GroupName).
				Validate(func(str string) error {
					return nil
				}),
			huh.NewSelect[int]().
				Title("Group Type").
				Value(&groupType).
				Validate(func(str int) error {
					createView.MemberGroup.GroupType = int64(str)
					return nil
				}).
				Options(groupSelectOptions...),

			huh.NewInput().
				Title("DN of LDAP group").
				Value(&createView.MemberGroup.GroupName).
				Validate(func(str string) error {
					return nil
				}),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}
