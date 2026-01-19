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
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

func TestLinkHandler(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"SameInputAsOutput", "harbor-artifact-tags.md", "harbor-artifact-tags.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := linkHandler(tt.input)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
func TestHasSeeAlso(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *cobra.Command
		expected bool
	}{
		{
			name: "Root command with no children",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use: "root",
				}
			},
			expected: false,
		},
		{
			name: "Child with a parent",
			setup: func() *cobra.Command {
				par := &cobra.Command{
					Use: "parent",
				}
				child := &cobra.Command{
					Use: "child",
				}
				par.AddCommand(child)
				return child
			},
			expected: true,
		},
		{
			name: "Root command with runnable child",
			setup: func() *cobra.Command {
				par := &cobra.Command{
					Use: "parent",
				}
				child := &cobra.Command{
					Use: "child",
					Run: func(cmd *cobra.Command, args []string) {},
				}
				par.AddCommand(child)
				return par
			},
			expected: true,
		},
		{
			name: "Root command with only a hidden child",
			setup: func() *cobra.Command {
				par := &cobra.Command{
					Use: "parent",
				}
				child := &cobra.Command{
					Use:    "child",
					Hidden: true,
				}
				par.AddCommand(child)
				return par
			},
			expected: false,
		},
		{
			name: "Root command with only a additional help topic child",
			setup: func() *cobra.Command {
				par := &cobra.Command{
					Use: "parent",
				}
				child := &cobra.Command{
					Use: "child",
				}
				par.AddCommand(child)
				return par
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasSeeAlso(tt.setup())
			if got != tt.expected {
				t.Errorf("hasSeeAlso() = %v, want %v", got, tt.expected)
			}
		})
	}
}
func TestGetWeight(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T) string
		expected     int
		expectedlogs string
	}{
		{
			name: "Valid file with weight",
			setup: func(t *testing.T) string {
				tmp := t.TempDir()
				filename := filepath.Join(tmp, "testfile")
				content := "---\ntitle: harbor artifact delete\nweight: 35\n---"
				if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
				return filename
			},
			expected:     35,
			expectedlogs: "",
		},
		{
			name: "No file path",
			setup: func(t *testing.T) string {
				return "non-existing-path"
			},
			expected:     0,
			expectedlogs: "unable to read file",
		},
		{
			name: "Only one dashed line",
			setup: func(t *testing.T) string {
				tmp := t.TempDir()
				filename := filepath.Join(tmp, "testfile")
				content := "---"
				if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
				return filename
			},
			expected:     0,
			expectedlogs: "YAML front matter not found on file",
		},
		{
			name: `Random content separated by \n`,
			setup: func(t *testing.T) string {
				tmp := t.TempDir()
				filename := filepath.Join(tmp, "testfile")
				content := "abcd\nefgh\nijkl\nmnop"
				if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
				return filename
			},
			expected:     0,
			expectedlogs: "YAML front matter not found on file",
		},
		{
			name: "Switched order of title and weight in the file",
			setup: func(t *testing.T) string {
				tmp := t.TempDir()
				filename := filepath.Join(tmp, "testfile")
				content := `---
weight: 20
title: test-title
---
`
				if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
				return filename
			},
			expected:     20,
			expectedlogs: "",
		},
		{
			name: "Malformed YAML (Broken weight)",
			setup: func(t *testing.T) string {
				tmp := t.TempDir()
				filename := filepath.Join(tmp, "testfile")
				content := `---
title: test-title
weight: [20]
---
`
				if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
				return filename
			},
			expected:     0,
			expectedlogs: "Failed to parse YAML in file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			originalLogOutput := log.StandardLogger().Out
			log.SetOutput(buf)
			defer log.SetOutput(originalLogOutput)

			filename := tt.setup(t)
			got := getWeight(filename)
			if got != tt.expected {
				t.Errorf("getWeight() = %v, want %v", got, tt.expected)
			}
			if !strings.Contains(buf.String(), tt.expectedlogs) {
				t.Errorf("Expected logs to contain %q, got logs:\n%s", tt.expectedlogs, buf.String())
			}
			// tempDir will get destroyed here
		})
	}
}
func TestPreblock(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*testing.T) string
		filebasename   string
		weight         int
		expectedtitle  string
		expectedweight int
		fileExists     bool
	}{
		{
			name: "Valid filename with file having valid title and weight",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				return dir
			},
			filebasename:   "harbor-test-command.md",
			weight:         35,
			expectedtitle:  "harbor test command",
			expectedweight: 35,
			fileExists:     true,
		},
		{
			name: "Non existing file",
			setup: func(t *testing.T) string {
				return "non-existing-dir"
			},
			filebasename:  "harbor-test-command.md",
			expectedtitle: "harbor test command",
			fileExists:    false,
		},
		{
			name: ".md extension not at end but is present in the directory name", // t.TempDir() creates a directory for testing using the name of the testcase, which in this case contains ".md"
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				return dir
			},
			filebasename:   "harbor-test-command",
			weight:         35,
			expectedtitle:  "harbor test command",
			expectedweight: 35,
			fileExists:     true,
		},
		{
			name: "Random filename",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				return dir
			},
			filebasename:   "randomfilename",
			weight:         35,
			expectedtitle:  "randomfilename",
			expectedweight: 35,
			fileExists:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tt.setup(t), tt.filebasename)
			if tt.fileExists {
				titleWithExtension := strings.ReplaceAll(tt.filebasename, "-", " ")
				title := strings.Split(titleWithExtension, ".md")[0]
				if err := os.WriteFile(path, []byte(fmt.Sprintf(frontmdtemplate, title, tt.weight)), 0600); err != nil {
					t.Fatal(err)
				}
			}
			pblock := preblock(path)
			var fm FrontMatter
			err := yaml.Unmarshal([]byte(pblock), &fm)
			if err != nil {
				t.Fatal("preblock returned invalid YAML", err)
			}
			if fm.Title != tt.expectedtitle {
				t.Errorf("Expected title %q, got %q", tt.expectedtitle, fm.Title)
			}
			if tt.fileExists {
				if fm.Weight != tt.expectedweight {
					t.Errorf("Expected weight %d, got %d", tt.expectedweight, fm.Weight)
				}
			} else {
				if fm.Weight < 1 {
					t.Errorf("Expected generated weight > 0, got %d", fm.Weight)
				}
				if fm.Weight%5 != 0 {
					t.Errorf("Expected weight to be multiple of 5, got %d", fm.Weight)
				}
			}
		})
	}
}
func TestPrintOptions(t *testing.T) {
	optionsPattern := func(flags, usage string) string {
		rgx := `### Options\s+` + "```sh" + `[\s\S]*?` + `\s+%s\s+.*%s` + `[\s\S]*?` + "```"
		return fmt.Sprintf(rgx, flags, usage)
	}
	optionsInheritedFromParentPattern := func(flags, usage string) string {
		rgx := `### Options inherited from parent commands\s+` + "```sh" + `[\s\S]*?` + `\s+%s\s+.*%s` + `[\s\S]*?` + "```"
		return fmt.Sprintf(rgx, flags, usage)
	}
	tests := []struct {
		name                 string
		setup                func() *cobra.Command
		expected             []string
		expectedNotToContain []string
	}{
		{
			name: "Command with only non-inherited flags",
			setup: func() *cobra.Command {
				cmd := &cobra.Command{
					Use: "test",
				}
				cmd.Flags().StringP("name", "n", "", "name flag")
				cmd.Flags().BoolP("verbose", "v", false, "verbose flag")
				return cmd
			},
			expected: []string{
				optionsPattern("-n, --name", "name flag"),
				optionsPattern("-v, --verbose", "verbose flag"),
			},
			expectedNotToContain: []string{
				optionsInheritedFromParentPattern("", ""),
			},
		},
		{
			name: "Command with inherited flags from parent",
			setup: func() *cobra.Command {
				parent := &cobra.Command{
					Use: "parent",
				}
				parent.PersistentFlags().StringP("config", "c", "", "config file")
				parent.PersistentFlags().BoolP("debug", "d", false, "debug mode")

				child := &cobra.Command{
					Use: "child",
				}
				child.Flags().StringP("output", "o", "", "output something")
				parent.AddCommand(child)
				return child
			},
			expected: []string{
				optionsPattern("-o, --output", "output something"),
				optionsInheritedFromParentPattern("-c, --config", "config file"),
				optionsInheritedFromParentPattern("-d, --debug", "debug mode"),
			},
		},
		{
			name: "Command with no flags",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use: "empty",
				}
			},
			expectedNotToContain: []string{
				optionsPattern("", ""),
				optionsInheritedFromParentPattern("", ""),
			},
		},
		{
			name: "Command with only inherited flags",
			setup: func() *cobra.Command {
				parent := &cobra.Command{
					Use: "parent",
				}
				parent.PersistentFlags().StringP("config", "c", "", "config file")
				child := &cobra.Command{
					Use: "child",
				}
				parent.AddCommand(child)
				return child
			},
			expected: []string{
				optionsInheritedFromParentPattern("-c, --config", "config file"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := tt.setup()

			err := printOptions(&buf, cmd)
			if err != nil {
				t.Fatalf("Error occurred in printOptions(): %v", err)
			}
			output := buf.String()
			for _, e := range tt.expected {
				rgx := regexp.MustCompile(e)
				if !rgx.MatchString(output) {
					t.Errorf("Expected output to contain the regex %q, but got:\n%s", e, output)
				}
			}
			for _, e := range tt.expectedNotToContain {
				rgx := regexp.MustCompile(e)
				if rgx.MatchString(output) {
					t.Errorf("Expected output NOT to contain the regex %q, but got:\n%s", e, output)
				}
			}
		})
	}
}
func TestMarkdownCustom(t *testing.T) {
	namePattern := func(name string) string {
		rgx := `##\s%s`
		return fmt.Sprintf(rgx, name)
	}
	descriptionPattern := func(short string) string {
		rgx := `### Description\s+#####\s+%s`
		return fmt.Sprintf(rgx, short)
	}
	longDescriptionPattern := func(long string) string {
		rgx := `### Synopsis\s+%s`
		return fmt.Sprintf(rgx, long)
	}
	examplePattern := func(example string) string {
		rgx := `### Examples\s+` + "```sh" + `\s+%s`
		return fmt.Sprintf(rgx, example)
	}
	uselinePattern := func(useline string) string {
		rgx := `%s \[flags\]`
		return fmt.Sprintf(rgx, useline)
	}
	optionsPattern := func(flags, usage string) string {
		rgx := `### Options\s+` + "```sh" + `[\s\S]*?` + `\s+%s\s+.*%s` + `[\s\S]*?` + "```"
		return fmt.Sprintf(rgx, flags, usage)
	}
	optionsInheritedFromParentPattern := func(flags, usage string) string {
		rgx := `### Options inherited from parent commands\s+` + "```sh" + `[\s\S]*?` + `\s+%s\s+.*%s` + `[\s\S]*?` + "```"
		return fmt.Sprintf(rgx, flags, usage)
	}
	seeAlsoPattern := func(name, link string) string {
		rgx := `### SEE ALSO[\s\S]*?` + `\*` + `\s+` + `\[%s\]\(%s\)`
		return fmt.Sprintf(rgx, name, link)
	}

	tests := []struct {
		name        string
		setup       func() *cobra.Command
		expected    []string
		notExpected []string
	}{
		{
			name: "Command with no parent or children",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use:   "testcmd",
					Short: "test command",
				}
			},
			expected: []string{
				namePattern("testcmd"),
				descriptionPattern("test command"),
				optionsPattern("-h, --help", "help for testcmd"), // "Options" is expected to be present because of the presence of cmd.InitDefaultHelpCmd() in MarkdownCustom function
			},
			notExpected: []string{
				longDescriptionPattern(""),
				examplePattern(""),
				optionsInheritedFromParentPattern("", ""),
				seeAlsoPattern("", ""),
			},
		},
		{
			name: "Command with no short description",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use: "testcmd",
				}
			},
			notExpected: []string{
				descriptionPattern(""),
			},
		},
		{
			name: "Command with short and long description",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use:   "testcmd",
					Short: "test command",
					Long:  "This is a long description\nwith multiple lines\nfor the test command",
				}
			},
			expected: []string{
				descriptionPattern("test command"),
				longDescriptionPattern("This is a long description\nwith multiple lines\nfor the test command"),
			},
		},
		{
			name: "Command with examples",
			setup: func() *cobra.Command {
				return &cobra.Command{
					Use:     "testcmd",
					Short:   "test command",
					Example: "test run myapp\ntest run otherapp",
				}
			},
			expected: []string{
				examplePattern("test run myapp\ntest run otherapp"),
			},
		},
		{
			name: "Command with parent",
			setup: func() *cobra.Command {
				parent := &cobra.Command{
					Use:   "parent",
					Short: "parent command",
				}
				child := &cobra.Command{
					Use:   "child",
					Short: "child command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				child.Flags().StringP("name", "n", "", "name flag")
				parent.AddCommand(child)
				return child
			},
			expected: []string{
				namePattern("parent child"),
				descriptionPattern("child command"),
				uselinePattern("parent child"), // this is the useline which is send to buffer if the command is runnable
				optionsPattern("-h, --help", "help for child"),
				optionsPattern("-n, --name", "name flag"),
				seeAlsoPattern("parent", "parent.md"),
			},
			notExpected: []string{
				optionsInheritedFromParentPattern("", ""),
			},
		},
		{
			name: "Command with children",
			setup: func() *cobra.Command {
				parent := &cobra.Command{
					Use:   "parent",
					Short: "parent command",
				}
				child1 := &cobra.Command{
					Use:   "child1",
					Short: "child1 command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				child2 := &cobra.Command{
					Use:   "child2",
					Short: "child2 command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				parent.AddCommand(child1, child2)
				return parent
			},
			expected: []string{
				namePattern("parent"),
				seeAlsoPattern("parent child1", "parent-child1.md"),
				seeAlsoPattern("parent child2", "parent-child2.md"),
			},
		},
		{
			name: "Command with flags",
			setup: func() *cobra.Command {
				cmd := &cobra.Command{
					Use:   "test",
					Short: "test command",
				}
				cmd.Flags().StringP("name", "n", "", "name flag")
				cmd.Flags().BoolP("verbose", "v", false, "verbose flag")
				return cmd
			},
			expected: []string{
				optionsPattern("-n, --name", "name flag"),
				optionsPattern("-v, --verbose", "verbose flag"),
			},
			notExpected: []string{
				optionsInheritedFromParentPattern("", ""),
				seeAlsoPattern("", ""),
			},
		},
		{
			name: "Command with hidden child (should not appear in SEE ALSO)",
			setup: func() *cobra.Command {
				parent := &cobra.Command{
					Use:   "parent",
					Short: "parent command",
				}
				visibleChild := &cobra.Command{
					Use:   "visible",
					Short: "visible child",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				hiddenChild := &cobra.Command{
					Use:    "hidden",
					Short:  "hidden child",
					Hidden: true,
				}
				parent.AddCommand(visibleChild, hiddenChild)
				return parent
			},
			expected: []string{
				seeAlsoPattern("parent visible", "parent-visible.md"),
			},
			notExpected: []string{
				seeAlsoPattern("parent hidden", "parent-hidden.md"),
			},
		},
		{
			name: "Command with all features",
			setup: func() *cobra.Command {
				gparent := &cobra.Command{
					Use:   "gparent",
					Short: "gparent command",
				}
				parent := &cobra.Command{
					Use:     "parent",
					Short:   "parent command",
					Long:    "this is a long description\n",
					Example: "harbor parent list\nharbor parent delete",
					Run:     func(cmd *cobra.Command, args []string) {},
				}
				parent.Flags().StringP("repository", "r", "", "repository name")
				child := &cobra.Command{
					Use:   "child",
					Short: "child command",
					Run:   func(parent *cobra.Command, args []string) {},
				}
				gparent.PersistentFlags().StringP("abcd", "a", "", "persistent flag")
				gparent.Flags().StringP("bcde", "b", "", "grandparent flag")
				parent.Flags().StringP("name", "n", "", "name flag")
				parent.AddCommand(child)
				gparent.AddCommand(parent)
				return parent
			},
			expected: []string{
				namePattern("gparent parent"),
				descriptionPattern("parent command"),
				longDescriptionPattern("this is a long description\n"),
				uselinePattern("gparent parent"), // useline for runnable commands
				examplePattern("harbor parent list\nharbor parent delete"),
				optionsPattern("-n, --name", "name flag"),
				optionsInheritedFromParentPattern("-a, --abcd", "persistent flag"),
				seeAlsoPattern("gparent", "gparent.md"),
				seeAlsoPattern("gparent parent child", "gparent-parent-child.md"),
			},
			notExpected: []string{
				optionsInheritedFromParentPattern("-b, --bcde", "grandparent flag"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			cmd := tt.setup()

			err := MarkdownCustom(cmd, &buf, linkHandler)
			if err != nil {
				t.Fatalf("MarkdownCustom() returned error: %v", err)
			}

			output := buf.String()

			for _, e := range tt.expected {
				rgx := regexp.MustCompile(e)
				if !rgx.MatchString(output) {
					t.Errorf("Expected output to contain the regex %q, but got:\n%s", e, output)
				}
			}

			for _, e := range tt.notExpected {
				rgx := regexp.MustCompile(e)
				if rgx.MatchString(output) {
					t.Errorf("Expected output NOT to contain the regex %q, but got:\n%s", e, output)
				}
			}
		})
	}
}
