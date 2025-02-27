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
	"os"
	"time"

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

	t := time.Now()

	header := &doc.GenManHeader{
		Title:   "HARBOR",
		Section: "1",
		Source:  "Habor Community",
		Manual:  "Harbor User Mannuals",
		Date:    &t,
	}

	err = doc.GenManTree(cmd.RootCmd(), header, docDir)
	if err != nil {
		fmt.Println("Error generating documentation:", err)
		os.Exit(1)
	}

	fmt.Println("Documentation generated successfully in", docDir)
	return nil
}

func main() {
	err := ManDoc()
	if err != nil {
		log.Fatal(err)
	}
}
