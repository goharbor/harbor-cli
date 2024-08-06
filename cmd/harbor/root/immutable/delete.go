package immutable

// import "github.com/spf13/cobra"

// func DeleteProjectCommand() *cobra.Command {

// 	cmd := &cobra.Command{
// 		Use:   "delete",
// 		Short: "delete immutable rule",
// 		Args:  cobra.MaximumNArgs(1),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			var err error

// 			if len(args) > 0 {
// 				err = api.DeleteProject(args[0])
// 			} else {
// 				projectName := prompt.GetProjectNameFromUser()
// 				err = api.DeleteProject(projectName)
// 			}
// 			if err != nil {
// 				log.Errorf("failed to delete project: %v", err)
// 			}
// 		},
// 	}

// 	return cmd
// }