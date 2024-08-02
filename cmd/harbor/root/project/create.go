package project

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/project/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CreateProjectCommand creates a new `harbor create project` command
func CreateProjectCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create [project name]",
		Short: "create project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			opts.ProjectName = args[0]
			var err error
			createView := &create.CreateView{
				ProjectName:  opts.ProjectName,
				Public:       opts.Public,
				RegistryID:   opts.RegistryID,
				StorageLimit: opts.StorageLimit,
				ProxyCache:   false,
			}
			if len(args) > 0 {
				opts.ProjectName = args[0]
				err = api.CreateProject(opts)
			} else if opts.ProjectName != "" && opts.RegistryID != "" && opts.StorageLimit != "" {
				err = api.CreateProject(opts)
			} else {
				err = createProjectView(createView)
			}

			if err != nil {
				log.Errorf("failed to create project: %v", err)
			}

		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.Public, "public", "", false, "Project is public or private")
	flags.StringVarP(&opts.RegistryID, "registry-id", "", "", "ID of referenced registry when creating the proxy cache project")
	flags.StringVarP(&opts.StorageLimit, "storage-limit", "", "-1", "Storage quota of the project")
	flags.BoolVarP(&opts.ProxyCache, "proxy-cache", "", false, "Whether the project is a proxy cache project")

	return cmd
}

func createProjectView(createView *create.CreateView) error {
	if createView == nil {
		createView = &create.CreateView{
			ProjectName:  "",
			Public:       false,
			RegistryID:   "",
			StorageLimit: "-1",
		}
	}

	create.CreateProjectView(createView)

	return api.CreateProject(*createView)

}
