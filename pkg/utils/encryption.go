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
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/zalando/go-keyring"
)

type KeyringProvider interface {
	Set(service, user, password string) error
	Get(service, user string) (string, error)
	Delete(service, user string) error
}

var keyringProvider KeyringProvider
var keyringProviderOnce sync.Once

func ensureKeyringProvider() {
	keyringProviderOnce.Do(func() {
		keyringProvider = GetKeyringProvider()
	})
}

type SystemKeyring struct{}

func (s *SystemKeyring) Set(service, user, password string) error {
	return keyring.Set(service, user, password)
}

func (s *SystemKeyring) Get(service, user string) (string, error) {
	return keyring.Get(service, user)
}

func (s *SystemKeyring) Delete(service, user string) error {
	return keyring.Delete(service, user)
}

// FileKeyring implements KeyringProvider using files in a directory
type FileKeyring struct {
	BaseDir string
}

func (f *FileKeyring) Set(service, user, password string) error {
	if err := os.MkdirAll(f.BaseDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filename := filepath.Join(f.BaseDir, sanitizeFilename(service+"_"+user))
	return os.WriteFile(filename, []byte(password), 0600)
}

func (f *FileKeyring) Get(service, user string) (string, error) {
	filename := filepath.Join(f.BaseDir, sanitizeFilename(service+"_"+user))
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *FileKeyring) Delete(service, user string) error {
	filename := filepath.Join(f.BaseDir, sanitizeFilename(service+"_"+user))
	return os.Remove(filename)
}

// Replace unsafe filename characters
func sanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			return '_'
		}
		return r
	}, name)
}

// EnvironmentKeyring implements KeyringProvider using environment variables
type EnvironmentKeyring struct {
	EnvVarName string
}

func (e *EnvironmentKeyring) Set(service, user, password string) error {
	// Set environment variable for the current process
	if err := os.Setenv(e.EnvVarName, password); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}
	return nil
}

func (e *EnvironmentKeyring) Get(service, user string) (string, error) {
	value := os.Getenv(e.EnvVarName)
	if value == "" {
		return "", fmt.Errorf("environment variable %s not found or empty", e.EnvVarName)
	}
	return value, nil
}

func (e *EnvironmentKeyring) Delete(service, user string) error {
	// Can't delete environment variables at runtime
	return fmt.Errorf("deleting environment variables at runtime is not supported")
}

// GetKeyringProvider selects the appropriate keyring provider
func GetKeyringProvider() KeyringProvider {
	// Priority 1: Check for environment variable configuration
	envKeyName := "HARBOR_ENCRYPTION_KEY"
	if envKey := os.Getenv(envKeyName); envKey != "" {
		logrus.Debug("Using environment-based encryption key")
		return &EnvironmentKeyring{
			EnvVarName: envKeyName,
		}
	}

	// Priority 2: Try system keyring
	if err := keyring.Set("harbor-cli-test", "test-user", "test"); err == nil {
		// Clean up the test entry
		err = keyring.Delete("harbor-cli-test", "test-user")
		if err != nil {
			logrus.Warnf("Failed to delete test entry from system keyring: %v", err)
		}
		logrus.Debug("Using system keyring")
		return &SystemKeyring{}
	}

	// Priority 3: Fall back to file-based keyring
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	fileKeyring := &FileKeyring{
		BaseDir: filepath.Join(homeDir, ".harbor", "keyring"),
	}

	logrus.Info("System keyring not available, using file-based keyring")
	return fileKeyring
}

func SetKeyringProvider(provider KeyringProvider) {
	keyringProvider = provider
}

const KeyringService = "harbor-cli"
const KeyringUser = "harbor-cli-encryption-key"

func GenerateEncryptionKey() error {
	ensureKeyringProvider()
	existingKey, err := keyringProvider.Get(KeyringService, KeyringUser)
	if err == nil && existingKey != "" {
		return nil
	}

	key := make([]byte, 32) // AES-256 key
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}
	return keyringProvider.Set(KeyringService, KeyringUser, base64.StdEncoding.EncodeToString(key))
}

func GetEncryptionKey() ([]byte, error) {
	ensureKeyringProvider()
	keyBase64, err := keyringProvider.Get(KeyringService, KeyringUser)
	if err != nil || keyBase64 == "" {
		// Attempt to generate a new key if not found
		if genErr := GenerateEncryptionKey(); genErr != nil {
			return nil, fmt.Errorf("failed to retrieve or generate encryption key: %w", err)
		}
		keyBase64, err = keyringProvider.Get(KeyringService, KeyringUser)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve encryption key after generation: %w", err)
		}
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	// Validate the key size for AES
	keySize := len(key)
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return nil, fmt.Errorf("invalid encryption key size: %d bytes. Must be 16, 24, or 32 bytes (after base64 decoding)", keySize)
	}

	return key, nil
}

func Encrypt(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key []byte, ciphertext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertextBytes := data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	return string(plaintext), nil
}
