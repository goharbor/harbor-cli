package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"gopkg.in/yaml.v2"
)

var configFile = filepath.Join(xdg.Home, ".harbor", "config")

type Credential struct {
	Name          string `yaml:"name"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	ServerAddress string `yaml:"serveraddress"`
}

type CredentialStore struct {
	CurrentCredentialName string       `yaml:"current-credential-name"`
	Credentials           []Credential `yaml:"credentials"`
}

func checkAndUpdateCredentialName(credential *Credential) {
	if credential.Name != "" {
		return
	}

	parsedUrl, err := url.Parse(credential.ServerAddress)
	if err != nil {
		panic(err)
	}

	credential.Name = parsedUrl.Hostname() + "-" + credential.Username
	log.Println("credential name not specified, storing the credential with the name as:", credential.Name)
	return
}

func readCredentialStore() (CredentialStore, error) {
	configInfo, err := ioutil.ReadFile(configFile)
	if err != nil {
		return CredentialStore{}, err
	}

	var credentialStore CredentialStore
	if err := yaml.Unmarshal(configInfo, &credentialStore); err != nil {
		return CredentialStore{}, err
	}
	return credentialStore, nil
}

func checkAndCreateConfigFile() {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create the parent directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(configFile), os.ModePerm); err != nil {
			panic(err)
		}

		if _, err := os.Create(configFile); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
}

func StoreCredential(credential Credential, setAsCurrentCredential bool) error {
	checkAndUpdateCredentialName(&credential)
	checkAndCreateConfigFile()
	credentialStore, err := readCredentialStore()
	if err != nil {
		return fmt.Errorf("failed to read credential store: %s", err)
	}

	// Check and remove the credential with same username and serveraddress.
	removeIndex := -1
	for i, cred := range credentialStore.Credentials {
		if cred.Username == credential.Username && cred.ServerAddress == credential.ServerAddress {
			removeIndex = i
			break
		}
	}

	if removeIndex != -1 {
		credentialStore.Credentials = append(credentialStore.Credentials[:removeIndex], credentialStore.Credentials[removeIndex+1:]...)
	}

	credentialStore.Credentials = append(credentialStore.Credentials, credential)
	if setAsCurrentCredential {
		credentialStore.CurrentCredentialName = credential.Name
	}

	bytes, err := yaml.Marshal(credentialStore)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(configFile, bytes, 0600); err != nil {
		return err
	}
	log.Println("Saving credentials to:", configFile)
	return nil
}

// resolveCredential resolves the credential in the following priority order:
// 1. credentialName specified by the user via CLI argument
// 2. credentialName specified by the user via environment variable
// 3. current active credential
func resolveCredential(credentialName string) (Credential, error) {
	credentialStore, err := readCredentialStore()
	if err != nil {
		panic(fmt.Sprintf("failed to read credential store: %s", err))
	}

	// If credentialName is not specified, check environment variable
	if credentialName == "" {
		credentialName = os.Getenv(constants.HarborCredentialName)
	}

	// If user has not specified the credential to use, use the current active credential
	if credentialName == "" {
		credentialName = credentialStore.CurrentCredentialName
	}

	if credentialName == "" {
		return Credential{}, fmt.Errorf("current credential name not set, please login again")
	}

	// Look for the credential with the given name
	for _, cred := range credentialStore.Credentials {
		if cred.Name == credentialName {
			return cred, nil
		}
	}

	return Credential{}, fmt.Errorf("no credential found for the name: %s, please login again with the credential name: %s", credentialName, credentialName)
}
