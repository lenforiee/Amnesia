package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type UserConfig struct {
	ServerURI      string `json:"server_uri"`
	PrivateKeyPath string `json:"private_key_path"`
	RemeberMe      bool   `json:"remeber_me"`
}

func NewUserConfig() (*UserConfig, error) {
	tempDir := os.TempDir()

	config := UserConfig{
		ServerURI:      "",
		PrivateKeyPath: "",
		RemeberMe:      false,
	}

	file, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/passbolt-gui/config.json", tempDir), file, 0644)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadUserConfig() (*UserConfig, error) {
	tempDir := os.TempDir()

	_, err := os.Stat(fmt.Sprintf("%s/passbolt-gui/config.json", tempDir))
	if os.IsNotExist(err) {
		return NewUserConfig()
	}

	file, err := ioutil.ReadFile(fmt.Sprintf("%s/passbolt-gui/config.json", tempDir))
	if err != nil {
		return nil, err
	}

	var config UserConfig
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (conf *UserConfig) SaveUserConfig() error {
	tempDir := os.TempDir()

	file, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/passbolt-gui/config.json", tempDir), file, 0644)
	if err != nil {
		return err
	}

	return nil
}
