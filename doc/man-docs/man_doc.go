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
	"github.com/spf13/cobra/doc"
)

func ManDoc() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	folderName := "man-docs/man1"
	_, err = os.Stat(folderName)
	if os.IsNotExist(err) {
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			log.Fatal("Error creating folder:", err)
		}
	}
	docDir := fmt.Sprintf("%s/%s", currentDir, folderName)

	header := &doc.GenManHeader{
		Title:   "HARBOR",
		Section: "1",
		Source:  "Harbor Community",
		Manual:  "Harbor User Manuals",
	}

	err = doc.GenManTree(cmd.RootCmd(), header, docDir)
	if err != nil {
		fmt.Println("Error generating documentation:", err)
		os.Exit(1)
	}

	err = cleanManPages(docDir)
	if err != nil {
		log.Fatalf("Error cleaning up documentation: %v", err)
	}

	fmt.Println("Documentation generated successfully in", docDir)
	return nil
}

func cleanManPages(docDir string) error {
	root, err := os.OpenRoot(docDir)
	if err != nil {
		return err
	}
	defer root.Close()

	return filepath.WalkDir(docDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".1") {
			relPath, err := filepath.Rel(docDir, path)
			if err != nil {
				return err
			}
			f, err := root.Open(relPath)
			if err != nil {
				return err
			}
			content, err := io.ReadAll(f)
			f.Close()
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

			wf, err := root.OpenFile(relPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			_, err = wf.Write([]byte(cleanedContent))
			wf.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func main() {
	err := ManDoc()
	if err != nil {
		log.Fatal(err)
	}
}
