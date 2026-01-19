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
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	cmd "github.com/goharbor/harbor-cli/cmd/harbor/root"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

const (
	markdownExtension = ".md"
	frontmdtemplate   = `---
title: %s
weight: %d
---
`
)

type FrontMatter struct {
	Title  string `yaml:"title"`
	Weight int    `yaml:"weight"`
}

func Doc() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	folderName := "cli-docs"
	_, err = os.Stat(folderName)
	if os.IsNotExist(err) {
		log.Printf("Folder %s does not exist", folderName)
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			log.Printf("Failed to create directory %s : %v", folderName, err)
			log.Fatal("Error creating folder:", err)
		}
	}
	docDir := fmt.Sprintf("%s/%s", currentDir, folderName)
	err = MarkdownTreeCustom(cmd.RootCmd(), docDir, preblock, linkHandler)
	if err != nil {
		return err
	}

	fmt.Println("Documentation generated at " + docDir)
	return nil
}

func preblock(filename string) string {
	randomNumber := rand.Intn(19) + 1 //nolint:gosec
	weight := randomNumber * 5

	prevWeight := getWeight(filename)
	if prevWeight > 0 {
		weight = prevWeight
	}

	baseName := filepath.Base(filename)
	name := strings.TrimSuffix(baseName, markdownExtension)
	title := strings.ReplaceAll(name, "-", " ")

	return fmt.Sprintf(frontmdtemplate, title, weight)
}

func linkHandler(s string) string {
	return s
}

func printOptions(buf *bytes.Buffer, cmd *cobra.Command) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("### Options\n\n```sh\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("### Options inherited from parent commands\n\n```sh\n")
		parentFlags.PrintDefaults()
		buf.WriteString("```\n\n")
	}
	return nil
}

func MarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("## " + name + "\n\n")
	if len(cmd.Short) > 0 {
		buf.WriteString("### Description\n\n")
		buf.WriteString("##### " + cmd.Short + "\n\n")
	}
	if len(cmd.Long) > 0 {
		buf.WriteString("### Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```sh\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(fmt.Sprintf("```sh\n%s\n```\n\n", cmd.Example))
	}

	if err := printOptions(buf, cmd); err != nil {
		return err
	}
	if hasSeeAlso(cmd) {
		buf.WriteString("### SEE ALSO\n\n")
		if cmd.HasParent() {
			parent := cmd.Parent()
			pname := parent.CommandPath()
			link := pname + markdownExtension
			link = strings.ReplaceAll(link, " ", "-")
			buf.WriteString(fmt.Sprintf("* [%s](%s)\t - %s\n", pname, linkHandler(link), parent.Short))
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Commands()

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			link := cname + markdownExtension
			link = strings.ReplaceAll(link, " ", "-")
			buf.WriteString(fmt.Sprintf("* [%s](%s)\t - %s\n", cname, linkHandler(link), child.Short))
		}
		buf.WriteString("\n")
	}

	_, err := buf.WriteTo(w)
	return err
}

func MarkdownTreeCustom(cmd *cobra.Command, dir string, filePrepender, linkHandler func(string) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := MarkdownTreeCustom(c, dir, filePrepender, linkHandler); err != nil {
			return err
		}
	}

	basename := strings.ReplaceAll(cmd.CommandPath(), " ", "-") + markdownExtension
	filename := filepath.Join(dir, basename)

	// if the file doesn't exist create or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	preblock := filePrepender(filename)
	f.Close()

	// create a fresh file and write the docs
	f, err = os.Create(filename)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(f, preblock); err != nil {
		return err
	}
	defer f.Close()
	if err := MarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}
	return nil
}

func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

// get previous weights from the yaml frontmatter
func getWeight(filename string) int {
	// Read the entire file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Warningf("unable to read file: %v", filename)
		return 0
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Check and extract YAML front matter
	if len(lines) < 3 || lines[0] != "---" {
		log.Warningf("YAML front matter not found on file: %v", filename)
		return 0
	}

	var yamlLines []string
	for _, line := range lines[1:] {
		if line == "---" {
			break
		}
		yamlLines = append(yamlLines, line)
	}

	yamlContent := strings.Join(yamlLines, "\n")

	// Unmarshal into Go struct
	var fm FrontMatter
	err = yaml.Unmarshal([]byte(yamlContent), &fm)
	if err != nil {
		log.Warningf("Failed to parse YAML in file %s", filename)
		return 0
	}
	return fm.Weight
}

func main() {
	err := Doc()
	if err != nil {
		log.Fatal(err)
	}
}
