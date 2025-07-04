package root

import (
	"fmt"
	"time"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/logs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Logs() *cobra.Command {
	var opts api.ListFlags
	var follow bool

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get recent logs of the projects which the user is a member of",
		Run: func(cmd *cobra.Command, args []string) {
			FormatFlag := viper.GetString("output-format")

			if follow {
				// Follow mode - continuous tailing
				followLogs(opts, FormatFlag)
			} else {
				// Single fetch mode
				logs, err := api.AuditLogs(opts)
				if err != nil {
					log.Fatalf("failed to retrieve audit logs: %v", err)
				}

				if FormatFlag != "" {
					utils.PrintPayloadInJSONFormat(logs.Payload)
					return
				}
				list.ListLogs(logs.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(
		&opts.Sort,
		"sort",
		"",
		"",
		"Sort the resource list in ascending or descending order",
	)
	flags.BoolVarP(&follow, "follow", "f", false, "Follow log output (tail -f behavior)")

	return cmd
}

func followLogs(opts api.ListFlags, formatFlag string) {
	var lastLogTime *time.Time

	// Set up logrus for clean output
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})

	fmt.Println("Following Harbor audit logs... (Press Ctrl+C to stop)")
	fmt.Println()

	for {
		logs, err := api.AuditLogs(opts)
		if err != nil {
			log.Errorf("failed to retrieve audit logs: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Filter new logs if we have a timestamp from the last fetch
		var newLogs []*models.AuditLogExt
		if lastLogTime != nil {
			for _, logEntry := range logs.Payload {
				// Convert strfmt.DateTime to time.Time
				logTime := time.Time(logEntry.OpTime)
				if !logTime.IsZero() && logTime.After(*lastLogTime) {
					newLogs = append(newLogs, logEntry)
				}
			}
		} else {
			// First run, show all logs
			newLogs = logs.Payload
		}

		// Update last log time with the most recent log
		if len(logs.Payload) > 0 {
			logTime := time.Time(logs.Payload[0].OpTime)
			if !logTime.IsZero() {
				lastLogTime = &logTime
			}
		}

		// Display new logs in streaming fashion
		if len(newLogs) > 0 {
			if formatFlag != "" {
				utils.PrintPayloadInJSONFormat(newLogs)
			} else {
				printLogsAsStream(newLogs)
			}
		}

		// Wait before next poll
		time.Sleep(2 * time.Second)
	}
}

func printLogsAsStream(logs []*models.AuditLogExt) {
	for _, logEntry := range logs {
		// Format the timestamp
		logTime := time.Time(logEntry.OpTime)
		timeStr := logTime.Format("2006-01-02 15:04:05")

		// Determine log level based on operation
		level := getLogLevel(logEntry.Operation)

		// Create a structured log entry
		// var logEntry *models.AuditLog
		entry := log.WithFields(log.Fields{
			"time":      timeStr,
			"user":      getUsername(logEntry.Username),
			"resource":  getResourceInfo(logEntry.ResourceType, logEntry.Resource),
			"operation": logEntry.Operation,
		})

		// Print with appropriate log level
		switch level {
		case "error":
			entry.Error(fmt.Sprintf("%s performed %s", getUsername(logEntry.Username), logEntry.Operation))
		case "warn":
			entry.Warn(fmt.Sprintf("%s performed %s", getUsername(logEntry.Username), logEntry.Operation))
		case "info":
			entry.Info(fmt.Sprintf("%s performed %s", getUsername(logEntry.Username), logEntry.Operation))
		default:
			entry.Debug(fmt.Sprintf("%s performed %s", getUsername(logEntry.Username), logEntry.Operation))
		}
	}
}

func getLogLevel(operation string) string {
	switch operation {
	case "delete", "stop", "remove":
		return "error"
	case "create", "push", "pull":
		return "info"
	case "update", "modify":
		return "warn"
	default:
		return "info"
	}
}

func getUsername(username string) string {
	if username == "" {
		return "unknown"
	}
	return username
}

func getResourceInfo(resourceType, resource string) string {
	if resourceType == "" && resource == "" {
		return "unknown"
	}
	if resourceType != "" && resource != "" {
		return fmt.Sprintf("%s:%s", resourceType, resource)
	}
	if resourceType != "" {
		return resourceType
	}
	return resource
}
