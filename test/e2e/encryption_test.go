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
