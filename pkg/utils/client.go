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
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/harbor"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	log "github.com/sirupsen/logrus"
)

var (
	ClientInstance *v2client.HarborAPI
	ClientOnce     sync.Once
	ClientErr      error
)

func GetClient() (*v2client.HarborAPI, error) {
	ClientOnce.Do(func() {
		config, err := GetCurrentHarborConfig()
		if err != nil {
			ClientErr = fmt.Errorf("failed to get current credential name: %v", err)
			return
		}
		credentialName := config.CurrentCredentialName
		if credentialName == "" {
			ClientErr = fmt.Errorf("no Harbor credentials found. Please run `harbor login` to configure access")
			return
		}

		ClientInstance, ClientErr = GetClientByCredentialName(credentialName)
		if ClientErr != nil {
			log.Errorf("failed to initialize client: %v", ClientErr)
			return
		}
	})

	return ClientInstance, ClientErr
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
		return nil, fmt.Errorf("failed to get credential %s: %w", credentialName, err)
	}
	if credential.AuthType == AuthTypeOIDC {
		return getOIDCClient(credential)
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

func getOIDCClient(credential Credential) (*v2client.HarborAPI, error) {
	idToken, err := GetDecryptedIDToken(credential.Name)
	if err != nil {
		return nil, err
	}

	if credential.ExpiresAt > 0 && time.Now().Unix() >= credential.ExpiresAt-60 {
		return nil, fmt.Errorf("OIDC session expired or is about to expire. Please run `harbor login %s --oidc` again", credential.ServerAddress)
	}

	return buildClientWithToken(credential.ServerAddress, idToken)
}

func buildClientWithToken(serverAddress, idToken string) (*v2client.HarborAPI, error) {
	u, err := url.Parse(serverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server URL: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid server URL: %s", serverAddress)
	}

	cfg := &harbor.Config{
		URL: u,
		AuthInfo: runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
			return req.SetHeaderParam("Authorization", "Bearer "+idToken)
		}),
	}

	return v2client.New(cfg.ToV2Config()), nil
}
