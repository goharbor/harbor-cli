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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Credential struct {
	Name          string `yaml:"name"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	ServerAddress string `yaml:"serveraddress"`
}

type HarborConfig struct {
	CurrentCredentialName string       `yaml:"current-credential-name"`
	Credentials           []Credential `yaml:"credentials"`
}

var (
	HarborFolder      string
	DefaultConfigPath string
)

func SetLocation() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	HarborFolder = filepath.Join(home, ".harbor")
	DefaultConfigPath = filepath.Join(HarborFolder, "config.yaml")
}

func (hc *HarborConfig) GetCurrentCredentialName() string {
	return hc.CurrentCredentialName
}

func CreateConfigFile() error {
	if _, err := os.Stat(DefaultConfigPath); os.IsNotExist(err) {
		if _, err := os.Create(DefaultConfigPath); err != nil {
			return err
		}
	}
	return nil
}

func AddCredentialsToConfigFile(credential Credential, configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	c := HarborConfig{}
	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	if c.Credentials == nil {
		c.Credentials = []Credential{}
	}

	c.Credentials = append(c.Credentials, credential)
	c.CurrentCredentialName = credential.Name

	viper.Set("current-credential-name", credential.Name)
	viper.Set("credentials", c.Credentials)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil

}

func GetCredentials(credentialName string) (Credential, error) {
	err := viper.ReadInConfig()
	if err != nil {
		return Credential{}, err
	}

	c := HarborConfig{}
	err = viper.Unmarshal(&c)
	if err != nil {
		return Credential{}, err
	}

	for _, cred := range c.Credentials {
		if cred.Name == credentialName {
			return cred, nil
		}
	}
	return Credential{}, nil
}
