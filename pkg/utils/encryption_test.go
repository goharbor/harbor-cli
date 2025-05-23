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
package utils_test

import (
	"testing"

	"github.com/goharbor/harbor-cli/pkg/utils"
)

func Test_EncryptionWithMockKeyring(t *testing.T) {
	// Use mock keyring for tests
	mockKeyring := utils.NewMockKeyring()
	utils.SetKeyringProvider(mockKeyring)

	// Run tests
	err := utils.GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("failed to generate encryption key: %v", err)
	}

	key, err := utils.GetEncryptionKey()
	if err != nil {
		t.Fatalf("failed to get encryption key: %v", err)
	}

	plaintext := "my-secret"
	encrypted, err := utils.Encrypt(key, []byte(plaintext))
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	decrypted, err := utils.Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("expected %s but got %s", plaintext, decrypted)
	}
}

func Test_EncryptionWithFileKeyring(t *testing.T) {
	// Use file keyring for tests
	fileKeyring := utils.FileKeyring{
		BaseDir: t.TempDir(),
	}
	utils.SetKeyringProvider(&fileKeyring)

	// Run tests
	err := utils.GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("failed to generate encryption key: %v", err)
	}

	key, err := utils.GetEncryptionKey()
	if err != nil {
		t.Fatalf("failed to get encryption key: %v", err)
	}

	plaintext := "my-secret"
	encrypted, err := utils.Encrypt(key, []byte(plaintext))
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	decrypted, err := utils.Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("expected %s but got %s", plaintext, decrypted)
	}
}

func Test_EncryptionWithEnvironmentKeyring(t *testing.T) {
	// Use environment keyring for tests
	envKeyring := utils.EnvironmentKeyring{
		EnvVarName: "TEST_HARBOR_ENCRYPTION_KEY",
	}
	utils.SetKeyringProvider(&envKeyring)

	// Run tests
	err := utils.GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("failed to generate encryption key: %v", err)
	}

	key, err := utils.GetEncryptionKey()
	if err != nil {
		t.Fatalf("failed to get encryption key: %v", err)
	}

	plaintext := "my-secret"
	encrypted, err := utils.Encrypt(key, []byte(plaintext))
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	decrypted, err := utils.Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("expected %s but got %s", plaintext, decrypted)
	}
}
