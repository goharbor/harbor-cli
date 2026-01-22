package gc

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ViewGCScheduleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule",
		Short: "Display the GC schedule",
		Run: func(cmd *cobra.Command, args []string) {
			scheduleWrapper, err := api.GetGCSchedule()
			if err != nil {
				logrus.Fatalf("Failed to get GC schedule: %v", err)
			}

			// scheduleWrapper is of type *models.GCHistory
			// It has a Schedule field which is *models.ScheduleObj
			// But verify if GCHistory itself has CreationTime/UpdateTime (it does as per inspection)
			// But we want the schedule's info.

			if scheduleWrapper == nil || scheduleWrapper.Schedule == nil {
				fmt.Println("No GC schedule set.")
				return
			}

			s := scheduleWrapper.Schedule

			fmt.Printf("Schedule Type:     %s\n", s.Type)
			if s.Cron != "" {
				fmt.Printf("Cron Expression:   %s\n", s.Cron)
			}
			fmt.Printf("Next Execution:    %v\n", s.NextScheduledTime)

			// models.ScheduleObj does not have CreationTime/UpdateTime fields based on inspection.
			// These fields usually belong to the wrapper GCHistory or Schedule wrapper.
			// If GCHistory has these fields, we should print from there.

			fmt.Printf("Creation Time:     %v\n", scheduleWrapper.CreationTime)
			fmt.Printf("Update Time:       %v\n", scheduleWrapper.UpdateTime)
		},
	}
}
