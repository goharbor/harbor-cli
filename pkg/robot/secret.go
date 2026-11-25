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
package robot

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/goharbor/harbor-cli/pkg/utils"
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/sirupsen/logrus"
)

// RobotSecret is the JSON document written by --export-to-file.
type RobotSecret struct {
	Name         string `json:"name"`
	ExpiresAt    int64  `json:"expires_at"`
	CreationTime string `json:"creation_time"`
	Secret       string `json:"secret"`
}

// ExportSecretToFile writes a robot's credentials to <name>-secret.json.
func ExportSecretToFile(name, secret, creationTime string, expiresAt int64) error {
	document, err := json.MarshalIndent(RobotSecret{
		Name:         name,
		ExpiresAt:    expiresAt,
		CreationTime: creationTime,
		Secret:       secret,
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal secret to JSON: %v", err)
	}

	// Remove a pre-existing file so the 0600 mode applies and no symlink
	// or wider permissions from an earlier export survive.
	filename := fmt.Sprintf("%s-secret.json", name)
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to replace existing secret file: %v", err)
	}
	if err := os.WriteFile(filename, document, 0600); err != nil {
		return fmt.Errorf("failed to write secret to file: %v", err)
	}

	fmt.Printf("Secret saved to %s\n", filename)
	return nil
}

// ReadSecretFromStdin reads and validates a robot secret from stdin.
func ReadSecretFromStdin() (string, error) {
	secret, err := utils.GetSecretStdin("Enter your secret: ")
	if err != nil {
		return "", fmt.Errorf("error reading secret: %v", err)
	}

	if err := utils.ValidatePassword(secret); err != nil {
		return "", fmt.Errorf("invalid secret: %v", err)
	}
	return secret, nil
}

// DisplaySecret shows a robot's secret and copies it to the clipboard.
func DisplaySecret(name, secret string) error {
	if err := robotview.FormSecret(name, secret); err != nil {
		return fmt.Errorf("failed to show the robot secret: %v", err)
	}

	if err := clipboard.WriteAll(secret); err != nil {
		logrus.Errorf("failed to write to clipboard: %v", err)
		return nil
	}
	fmt.Println("secret copied to clipboard.")
	return nil
}
