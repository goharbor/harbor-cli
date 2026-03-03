// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package root

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/logs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logsLogger = log.New()

func Logs() *cobra.Command {
	var opts api.ListFlags
	var follow bool
	var refreshInterval string

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get recent logs of the projects which the user is a member of",
		Args:  cobra.NoArgs,
		Long: `Get recent logs of the projects which the user is a member of.
This command retrieves the audit logs for the projects the user is a member of. It supports pagination, sorting, and filtering through query parameters. The logs can be followed in real-time with the --follow flag, and the output can be formatted as JSON with the --output-format flag.

harbor-cli logs --page 1 --page-size 10 --query "operation=push" --sort "op_time:desc"

harbor-cli logs --follow --refresh-interval 2s

harbor-cli logs --output-format json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}
			if refreshInterval != "" && !follow {
				fmt.Println("The --refresh-interval flag is only applicable when using --follow. It will be ignored.")
			}

			if follow {
				var interval time.Duration = 5 * time.Second
				var err error
				if refreshInterval != "" {
					interval, err = time.ParseDuration(refreshInterval)
					if err != nil {
						return fmt.Errorf("invalid refresh interval: %w", err)
					}
					if interval < 500*time.Millisecond {
						return fmt.Errorf("refresh-interval must be at least 500ms (got: %v)", interval)
					}
				}
				followLogs(opts, interval)
				return nil
			}

			logs, err := api.AuditLogs(opts)
			if err != nil {
				return fmt.Errorf("failed to retrieve audit logs: %w", err)
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				log.WithField("output_format", formatFlag).Debug("Output format selected")
				err = utils.PrintFormat(logs.Payload, formatFlag)
				if err != nil {
					return err
				}
			} else {
				list.ListLogs(logs.Payload)
			}
			return nil
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
	flags.StringVarP(&refreshInterval, "refresh-interval", "n", "",
		"Interval to refresh logs when following (default: 5s)")

	return cmd
}

func followLogs(opts api.ListFlags, interval time.Duration) {
	var lastLogTime *time.Time

	logsLogger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})
	logsLogger.SetLevel(log.InfoLevel)
	logsLogger.SetOutput(os.Stdout)

	fmt.Println("Following Harbor audit logs... (Press Ctrl+C to stop)")

	for {
		logs, err := api.AuditLogs(opts)
		if err != nil {
			log.Errorf("failed to retrieve audit logs: %v", err)
			time.Sleep(interval)
			continue
		}

		var newLogs []*models.AuditLogExt
		if lastLogTime != nil {
			for _, logEntry := range logs.Payload {
				logTime := time.Time(logEntry.OpTime)
				if !logTime.IsZero() && logTime.After(*lastLogTime) {
					newLogs = append(newLogs, logEntry)
				}
			}
		} else {
			newLogs = logs.Payload
		}

		if len(logs.Payload) > 0 {
			logTime := time.Time(logs.Payload[0].OpTime)
			if !logTime.IsZero() {
				lastLogTime = &logTime
			}
		}

		printLogsAsStream(newLogs)
		time.Sleep(interval)
	}
}

func printLogsAsStream(logs []*models.AuditLogExt) {
	for _, logEntry := range logs {
		logTime := time.Time(logEntry.OpTime)
		level := getLogLevel(logEntry.OperationResult)

		displayUser := truncateUsername(logEntry.Username)
		resource := getResourceInfo(logEntry.ResourceType, logEntry.Resource)

		resultIcon := "✓"
		if !logEntry.OperationResult {
			resultIcon = "✗"
		}

		message := fmt.Sprintf("%s %s %s %s",
			displayUser,
			logEntry.Operation,
			resource,
			resultIcon)

		entry := logsLogger.WithTime(logTime)

		switch level {
		case "error":
			entry.Error(message)
		case "info":
			entry.Info(message)
		default:
			entry.Debug(message)
		}
	}
}

func truncateUsername(username string) string {
	if username == "" {
		return "unknown"
	}

	if len(username) > 30 {
		if parts := strings.Split(username, "+"); len(parts) > 1 {
			project := strings.TrimPrefix(parts[0], "robt_")
			return fmt.Sprintf("%s+robot", project)
		}
		return username[:27] + "..."
	}
	return username
}

func getLogLevel(operationResult bool) string {
	switch operationResult {
	case false:
		return "error"
	case true:
		return "info"
	default:
		return "error"
	}
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
