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
	"io"
	"time"

	"github.com/goharbor/harbor-cli/cmd/harbor/root/context"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/replication"

	"github.com/goharbor/harbor-cli/cmd/harbor/root/artifact"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/cve"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/instance"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/labels"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/member"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/quota"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/registry"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/repository"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/scan_all"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/scanner"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/schedule"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/tag"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/user"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/webhook"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	output  string
	cfgFile string
	verbose bool
)

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "harbor",
		Short:        "Official Harbor CLI",
		SilenceUsage: true,
		Long:         "Official Harbor CLI",
		Example: `
// Base command:
harbor

// Display help about the command:
harbor help
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Determine if --config was explicitly set
			userSpecifiedConfig := cmd.Flags().Changed("config")
			// Initialize configuration
			utils.InitConfig(cfgFile, userSpecifiedConfig)

			// Conditionally set the timestamp format only in verbose mode
			formatter := &logrus.TextFormatter{}

			if verbose {
				formatter.FullTimestamp = true
				formatter.TimestampFormat = time.RFC3339
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				logrus.SetOutput(io.Discard)
			}
			logrus.SetFormatter(formatter)

			return nil
		},
	}

	root.PersistentFlags().StringVarP(&output, "output-format", "o", "", "Output format. One of: json|yaml")
	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/harbor-cli/config.yaml)")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	err := viper.BindPFlag("output-format", root.PersistentFlags().Lookup("output-format"))
	if err != nil {
		fmt.Println(err.Error())
	}

	err = viper.BindPFlag("config", root.PersistentFlags().Lookup("config"))
	if err != nil {
		fmt.Println(err.Error())
	}

	root.AddGroup(&cobra.Group{ID: "core", Title: "Core:"})
	root.AddGroup(&cobra.Group{ID: "access", Title: "Access:"})
	root.AddGroup(&cobra.Group{ID: "system", Title: "System:"})
	root.AddGroup(&cobra.Group{ID: "utils", Title: "Utility:"})

	// Core
	cmd := InfoCommand()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	cmd = project.Project()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	cmd = repository.Repository()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	cmd = artifact.Artifact()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	cmd = tag.TagCommand()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	cmd = labels.Labels()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	cmd = quota.Quota()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	cmd = cve.CVEAllowlist()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	cmd = webhook.Webhook()
	cmd.GroupID = "core"
	root.AddCommand(cmd)

	// Access
	cmd = LoginCommand()
	cmd.GroupID = "access"
	root.AddCommand(cmd)

	cmd = user.User()
	cmd.GroupID = "access"
	root.AddCommand(cmd)

	cmd = member.Member()
	cmd.GroupID = "access"
	root.AddCommand(cmd)

	// System
	cmd = context.Context()
	cmd.GroupID = "system"
	root.AddCommand(cmd)

	cmd = HealthCommand()
	cmd.GroupID = "system"
	root.AddCommand(cmd)

	cmd = instance.Instance()
	cmd.GroupID = "system"
	root.AddCommand(cmd)

	cmd = registry.Registry()
	cmd.GroupID = "system"
	root.AddCommand(cmd)

	cmd = replication.Replication()
	cmd.GroupID = "system"
	root.AddCommand(cmd)

	cmd = scanner.Scanner()
	cmd.GroupID = "system"
	root.AddCommand(cmd)

	cmd = scan_all.ScanAll()
	cmd.GroupID = "system"
	root.AddCommand(cmd)

	cmd = schedule.Schedule()
	cmd.GroupID = "system"
	root.AddCommand(cmd)

	// Utils
	cmd = versionCommand()
	cmd.GroupID = "utils"
	root.AddCommand(cmd)

	return root
}
