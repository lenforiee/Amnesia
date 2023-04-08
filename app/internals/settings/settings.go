package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

type UserConfig struct {
	ServerURI      string `json:"server_uri"`
	PrivateKeyPath string `json:"private_key_path"`
	RememberMe     bool   `json:"remember_me"`
}

func GetOSSaveDir() (userDir string) {

	if runtime.GOOS == "windows" {
		userDir = os.Getenv("APPDATA")
	} else if runtime.GOOS == "linux" {
		userDir = os.Getenv("HOME") + "/.config"
	} else if runtime.GOOS == "darwin" {
		userDir = os.Getenv("HOME") + "/Library/Application Support"
	}
	return userDir
}

func NewUserConfig() (*UserConfig, error) {
	userDir := GetOSSaveDir()

	config := UserConfig{
		ServerURI:      "",
		PrivateKeyPath: "",
		RememberMe:     false,
	}

	_, err := os.Stat(fmt.Sprintf("%s/amnesia", userDir))
	if os.IsNotExist(err) {
		err = os.Mkdir(fmt.Sprintf("%s/amnesia", userDir), 0755)
		if err != nil {
			return nil, err
		}
	}

	file, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/amnesia/config.json", userDir), file, 0644)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadUserConfig() (*UserConfig, error) {
	userDir := GetOSSaveDir()

	_, err := os.Stat(fmt.Sprintf("%s/amnesia/config.json", userDir))
	if os.IsNotExist(err) {
		return NewUserConfig()
	}

	file, err := ioutil.ReadFile(fmt.Sprintf("%s/amnesia/config.json", userDir))
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

	userDir := GetOSSaveDir()
	file, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/amnesia/config.json", userDir), file, 0644)
	if err != nil {
		return err
	}

	return nil
}
