package root

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)
func GenerateDocs() *cobra.Command {

	cmd := &cobra.Command{ 
		Use: "gen-docs",
		Short: "Generate documentation for Harbor CLI",
		Long:   `Generate documentation for Harbor CLI`,
		Hidden: true,
		Run: func (cmd *cobra.Command, args []string) {
			path, err:= cmd.Flags().GetString("dir")
			if err != nil {
				log.Fatal(err)
			}
			if path == "" {
				path = "docs"
			}
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			err = doc.GenYamlTree(cmd.Root(), path)
			if err != nil {
				log.Fatal(err)
			}
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			err = replaceHomeDirWithEnvVar(path, homeDir, "$HOME")
			if err != nil {
				log.Fatal(err)
			}

			log.Infof("Documentation generated in %s", path)
		},
	}

	cmd.Flags().String("dir", "", "Directory to generate documentation in")
	return cmd

}

func replaceHomeDirWithEnvVar(path, homeDir, envVar string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			updatedContent := strings.ReplaceAll(string(content), homeDir, envVar)
			err = os.WriteFile(filePath, []byte(updatedContent), info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
}