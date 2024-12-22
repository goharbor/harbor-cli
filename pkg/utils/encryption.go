package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/zalando/go-keyring"
)

type KeyringProvider interface {
	Set(service, user, password string) error
	Get(service, user string) (string, error)
	Delete(service, user string) error
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

var keyringProvider KeyringProvider = &SystemKeyring{}

func SetKeyringProvider(provider KeyringProvider) {
	keyringProvider = provider
}

const KeyringService = "harbor-cli"
const KeyringUser = "harbor-cli-encryption-key"

func GenerateEncryptionKey() error {
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
	return base64.StdEncoding.DecodeString(keyBase64)
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
