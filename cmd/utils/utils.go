package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/goharbor/go-client/pkg/harbor"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"gopkg.in/yaml.v2"
)

var authFile = filepath.Join(xdg.Home, ".harbor", "config", "auth.yaml")

// A wrapper to hold auth data
type AuthDataWrapper struct {
	ServerAddress string
	Username      string
	Password      string
}

func readAuthData() *AuthDataWrapper {
	authInfo, err := ioutil.ReadFile(authFile)
	if err != nil {
		panic(err)
	}

	authData := &AuthDataWrapper{}
	err = yaml.Unmarshal(authInfo, &authData)
	if err != nil {
		panic(err)
	}
	return authData
}

// Saves the auth data to the auth file
func SaveAuthData(authData *AuthDataWrapper) {
	bytes, err := yaml.Marshal(authData)
	if err != nil {
		panic(err)
	}
	os.MkdirAll(filepath.Dir(authFile), os.ModePerm)
	if err = ioutil.WriteFile(authFile, bytes, 0600); err != nil {
		panic(err)
	}
	fmt.Println("Auth data saved to ", authFile)
}

// Returns Harbor v2 client
func GetClient(clientConfig *harbor.ClientSetConfig) *v2client.HarborAPI {
	if clientConfig == nil {
		clientConfig = getClientConfigFromAuthFile()
	}
	cs, err := harbor.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}
	return cs.V2()
}

func getClientConfigFromAuthFile() *harbor.ClientSetConfig {
	authData := readAuthData()
	return &harbor.ClientSetConfig{
		URL:      authData.ServerAddress,
		Username: authData.Username,
		Password: authData.Password,
	}
}

func PrintPayloadInJSONFormat(payload any) {
	if payload == nil {
		return
	}

	jsonStr, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonStr))
}
