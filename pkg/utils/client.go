// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"os"
	"sync"

	"github.com/goharbor/go-client/pkg/harbor"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	clientInstance *v2client.HarborAPI
	clientOnce     sync.Once
	clientErr      error
)

func GetClient() (*v2client.HarborAPI, error) {
	clientOnce.Do(func() {
		credentialName := viper.GetString("current-credential-name")
		clientInstance = GetClientByCredentialName(credentialName)
		if clientErr != nil {
			log.Errorf("failed to initialize client: %v", clientErr)
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
func GetClientByCredentialName(credentialName string) *v2client.HarborAPI {
	credential, err := GetCredentials(credentialName)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	clientConfig := &harbor.ClientSetConfig{
		URL:      credential.ServerAddress,
		Username: credential.Username,
		Password: credential.Password,
	}
	return GetClientByConfig(clientConfig)
}
