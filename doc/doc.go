package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"

	cmd "github.com/goharbor/harbor-cli/cmd/harbor/root"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const frontmdtemplate = `---
title: %s
weight: %d
---
`

func preblock(filename string) string {
	file := strings.Split(filename, ".md")
	name := filepath.Base(file[0])
	title := strings.ReplaceAll(name, "-", " ")
	randomNumber := rand.Intn(20)
	weight := randomNumber * 5

	return fmt.Sprintf(frontmdtemplate, title, weight)
}

func Doc() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	folderName := "CLIDoc"
	_, err = os.Stat(folderName)
	if os.IsNotExist(err) {
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			log.Fatal("Error creating folder:", err)
		}
	}
	docDir := fmt.Sprintf("%s/%s", currentDir, folderName)
	err = MarkdownTreeCustom(cmd.RootCmd(), docDir, preblock, linkHandler)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Documentation generated at " + docDir)
}

func linkHandler(name string) string {
	base := strings.TrimSuffix(name, path.Ext(name))
	words := strings.Split(base, "-")
	if len(words) <= 1 {
		return ""
	}
	if len(words) == 3 {
		return strings.ToLower(words[2])
	}
	if len(words) == 4 {
		return strings.ToLower(words[2]) + "/" + strings.ToLower(words[3])
	}
	return strings.ToLower(words[1])
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
	buf.WriteString("#### " + cmd.Short + "\n\n")
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

	basename := strings.ReplaceAll(cmd.CommandPath(), " ", "-") + ".md"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(filename)); err != nil {
		return err
	}
	if err := MarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}
	return nil
}

func main() {
	Doc()
}
