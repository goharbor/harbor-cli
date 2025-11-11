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
package view

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// DiffConfigurations displays configuration changes in AWS CDK style
func DiffConfigurations(upstreamConfigs, localConfigs map[string]interface{}) {
	// Collect all unique field names and sort them
	allFields := make(map[string]bool)
	for field := range upstreamConfigs {
		allFields[field] = true
	}
	for field := range localConfigs {
		allFields[field] = true
	}

	var sortedFields []string
	for field := range allFields {
		sortedFields = append(sortedFields, field)
	}
	sort.Strings(sortedFields)

	// Track changes by type
	var additions, modifications, deletions []string
	changeDetails := make(map[string][2]string)

	for _, field := range sortedFields {
		upstreamVal, hasUpstream := upstreamConfigs[field]
		localVal, hasLocal := localConfigs[field]

		upstreamStr := formatValuePlain(upstreamVal, hasUpstream)
		localStr := formatValuePlain(localVal, hasLocal)

		if !hasUpstream && hasLocal {
			additions = append(additions, field)
			changeDetails[field] = [2]string{"", localStr}
		} else if hasUpstream && !hasLocal {
			deletions = append(deletions, field)
			changeDetails[field] = [2]string{upstreamStr, ""}
		} else if upstreamStr != localStr {
			modifications = append(modifications, field)
			changeDetails[field] = [2]string{upstreamStr, localStr}
		}
	}

	totalChanges := len(additions) + len(modifications) + len(deletions)
	if totalChanges == 0 {
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
		fmt.Println(successStyle.Render("✓ No changes detected."))
		return
	}

	// Define styles
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	modifyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	deleteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	fieldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	oldValueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	newValueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)

	// Header
	fmt.Println(headerStyle.Render("Configuration changes to be applied:"))
	fmt.Println(infoStyle.Render("(For available configuration fields, see: https://github.com/goharbor/go-client/blob/main/pkg/sdk/v2.0/models/configurations.go)"))

	fmt.Println()

	// Summary
	summary := []string{}
	if len(additions) > 0 {
		summary = append(summary, addStyle.Render(fmt.Sprintf("[+] %d to add", len(additions))))
	}
	if len(modifications) > 0 {
		summary = append(summary, modifyStyle.Render(fmt.Sprintf("[~] %d to modify", len(modifications))))
	}
	if len(deletions) > 0 {
		summary = append(summary, deleteStyle.Render(fmt.Sprintf("[-] %d to remove", len(deletions))))
	}
	fmt.Println(strings.Join(summary, "  "))
	fmt.Println()

	// Print additions
	if len(additions) > 0 {
		for _, field := range additions {
			values := changeDetails[field]
			fmt.Printf("%s %s\n", addStyle.Render("[+]"), fieldStyle.Render(field))
			fmt.Printf("    %s %s\n", grayStyle.Render("└─"), newValueStyle.Render(values[1]))
			fmt.Println()
		}
	}

	// Print modifications
	if len(modifications) > 0 {
		for _, field := range modifications {
			values := changeDetails[field]
			fmt.Printf("%s %s\n", modifyStyle.Render("[~]"), fieldStyle.Render(field))
			fmt.Printf("    %s %s\n", oldValueStyle.Render("[-]"), oldValueStyle.Render(values[0]))
			fmt.Printf("    %s %s\n", newValueStyle.Render("[+]"), newValueStyle.Render(values[1]))
			fmt.Println()
		}
	}

	// Print deletions
	if len(deletions) > 0 {
		for _, field := range deletions {
			values := changeDetails[field]
			fmt.Printf("%s %s\n", deleteStyle.Render("[-]"), fieldStyle.Render(field))
			fmt.Printf("    %s %s\n", grayStyle.Render("└─"), oldValueStyle.Render(values[0]))
			fmt.Println()
		}
	}
}

// formatValuePlain converts a value to a plain string without styling
func formatValuePlain(val interface{}, exists bool) string {
	if !exists {
		return "(not set)"
	}

	if val == nil {
		return "(nil)"
	}

	switch v := val.(type) {
	case string:
		if v == "" {
			return "(empty)"
		}
		// Quote strings to make them clear
		if len(v) > 60 {
			return fmt.Sprintf(`"%s..."`, v[:57])
		}
		return fmt.Sprintf(`"%s"`, v)
	case bool:
		return fmt.Sprintf("%t", v)
	case int, int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	default:
		str := fmt.Sprintf("%v", v)
		if len(str) > 60 {
			return str[:57] + "..."
		}
		return str
	}
}
