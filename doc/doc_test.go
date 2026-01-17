package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
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
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
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
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
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
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
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
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
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
				if err := os.WriteFile(path, []byte(fmt.Sprintf(frontmdtemplate, title, tt.weight)), 0644); err != nil {
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
