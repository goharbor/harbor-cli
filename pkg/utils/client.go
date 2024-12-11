package utils

import (
	"context"
	"fmt"
	"sync"

	"github.com/goharbor/go-client/pkg/harbor"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	log "github.com/sirupsen/logrus"
)

var (
	clientInstance *v2client.HarborAPI
	clientOnce     sync.Once
	clientErr      error
)

func GetClient() (*v2client.HarborAPI, error) {
	clientOnce.Do(func() {
		config, err := GetCurrentHarborConfig()
		if err != nil {
			clientErr = fmt.Errorf("failed to get current credential name: %v", err)
			return
		}
		credentialName := config.CurrentCredentialName
		if credentialName == "" {
			clientErr = fmt.Errorf("current-credential-name is not set in config file")
			return
		}

		clientInstance, clientErr = GetClientByCredentialName(credentialName)
		if clientErr != nil {
			log.Errorf("failed to initialize client: %v", clientErr)
			return
		}
	})

	return clientInstance, clientErr
}

func ContextWithClient() (context.Context, *v2client.HarborAPI, error) {
	client, err := GetClient()
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()
	return ctx, client, nil
}

func GetClientByConfig(clientConfig *harbor.ClientSetConfig) *v2client.HarborAPI {
	cs, err := harbor.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}
	return cs.V2()
}

// Returns Harbor v2 client after resolving the credential name
func GetClientByCredentialName(credentialName string) (*v2client.HarborAPI, error) {
	credential, err := GetCredentials(credentialName)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	// Get encryption key
	key, err := GetEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption key: %w", err)
	}

	// Decrypt password
	decryptedPassword, err := Decrypt(key, string(credential.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	clientConfig := &harbor.ClientSetConfig{
		URL:      credential.ServerAddress,
		Username: credential.Username,
		Password: decryptedPassword,
	}
	return GetClientByConfig(clientConfig), nil
}
