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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cmd "github.com/goharbor/harbor-cli/cmd/harbor/root"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type ManTreeGenerator func(cmd *cobra.Command, header *doc.GenManHeader, dir string) error
type ManPageCleaner func(string) error

func ManDoc(w io.Writer, generator ManTreeGenerator, cleaner ManPageCleaner) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	folderName := "man-docs/man1"
	_, err = os.Stat(folderName)
	if os.IsNotExist(err) {
		err = os.MkdirAll(folderName, 0755) //cannot create nested directories using os.Mkdir
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	docDir := fmt.Sprintf("%s/%s", currentDir, folderName)

	header := &doc.GenManHeader{
		Title:   "HARBOR",
		Section: "1",
		Source:  "Harbor Community",
		Manual:  "Harbor User Manuals",
	}

	err = generator(cmd.RootCmd(), header, docDir)
	if err != nil {
		return fmt.Errorf("error generating documentation: %w", err)
	}

	err = cleaner(docDir)
	if err != nil {
		return fmt.Errorf("error cleaning documentation: %w", err)
	}

	fmt.Fprintf(w, "Documentation generated successfully in %s\n", docDir)
	return nil
}

func cleanManPages(docDir string) error {
	return filepath.Walk(docDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".1") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)
			lines := strings.Split(contentStr, "\n")
			for i, line := range lines {
				if strings.HasPrefix(line, ".TH ") {
					re := regexp.MustCompile(`"[^"]*"`)
					matches := re.FindAllString(line, -1)
					if len(matches) >= 5 {
						matches[2] = ``
						lines[i] = ".TH " + strings.Join(matches, " ")
					}
				}
			}
			updatedContent := strings.Join(lines, "\n")

			cleanedContent := strings.Split(updatedContent, "\n.SH HISTORY")[0]
			cleanedContent = strings.TrimRight(cleanedContent, "\n")

			err = os.WriteFile(path, []byte(cleanedContent), 0600)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func main() {
	err := ManDoc(os.Stdout, doc.GenManTree, cleanManPages)
	if err != nil {
		log.Fatal(err)
	}
}
