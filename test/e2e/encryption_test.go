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
package e2e

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
