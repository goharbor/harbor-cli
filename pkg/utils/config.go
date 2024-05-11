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
	dir := filepath.Dir(DefaultConfigPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

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
