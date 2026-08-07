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
	"bufio"
	"fmt"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/declarative"
	"github.com/spf13/cobra"
)

func applyCommand() *cobra.Command {
	var filename string
	var dryRun bool
	var skipConfirmation bool
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Reconcile Harbor with a declarative configuration",
		Long: `Compare a versioned configuration document with Harbor and apply the
required create and update operations. Apply is additive: resources omitted
from the document are left untouched, and extra live resources are not deleted.`,
		Example: `  harbor apply -f harbor.yaml --dry-run
  harbor apply -f harbor.yaml
  harbor apply -f environments/production
  harbor apply -f harbor.yaml --yes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if filename == "" {
				return fmt.Errorf("configuration file is required")
			}
			desired, err := declarative.ReadPath(filename)
			if err != nil {
				return err
			}
			backend, err := declarative.NewLiveBackend()
			if err != nil {
				return err
			}
			service := declarative.NewService(backend)
			plan, err := service.Plan(cmd.Context(), desired)
			if err != nil {
				return err
			}
			printPlan(cmd, plan)
			if !plan.HasChanges() {
				fmt.Fprintln(cmd.OutOrStdout(), "Harbor configuration is already converged.")
				return nil
			}
			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: %d change(s) would be applied.\n", plan.ChangeCount())
				return nil
			}
			if !skipConfirmation {
				confirmed, err := confirmApply(cmd)
				if err != nil {
					return err
				}
				if !confirmed {
					fmt.Fprintln(cmd.OutOrStdout(), "Apply cancelled.")
					return nil
				}
			}
			if err := service.Apply(cmd.Context(), desired, plan); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Applied %d change(s).\n", plan.ChangeCount())
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, "file", "f", "", "YAML/JSON configuration file or directory to apply")
	flags.BoolVar(&dryRun, "dry-run", false, "Show the reconciliation plan without changing Harbor")
	flags.BoolVarP(&skipConfirmation, "yes", "y", false, "Apply without an interactive confirmation")
	return cmd
}

func printPlan(cmd *cobra.Command, plan declarative.Plan) {
	fmt.Fprintln(cmd.OutOrStdout(), "Reconciliation plan:")
	for _, action := range plan.Actions {
		if action.Operation == declarative.OperationNoop && !verbose {
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", action)
	}
}

func confirmApply(cmd *cobra.Command) (bool, error) {
	fmt.Fprint(cmd.OutOrStdout(), "Apply these changes? (y/N): ")
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("read confirmation: %w", err)
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}
