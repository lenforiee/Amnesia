package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

type UserSettings struct {
	UserAgent      string `json:"user_agent"`
	ServerURI      string `json:"server_uri"`
	PrivateKeyPath string `json:"private_key_path"`
	RememberMe     bool   `json:"remember_me"`
}

func GetSaveDir() (userDir string) {

	if runtime.GOOS == "windows" {
		userDir = os.Getenv("APPDATA")
	} else if runtime.GOOS == "linux" {
		userDir = os.Getenv("HOME") + "/.config"
	} else if runtime.GOOS == "darwin" {
		userDir = os.Getenv("HOME") + "/Library/Application Support"
	}
	return userDir
}

func NewUserSettings() (s UserSettings, err error) {
	userDir := GetSaveDir()

	s = UserSettings{
		UserAgent:      "",
		ServerURI:      "",
		PrivateKeyPath: "",
		RememberMe:     false,
	}

	_, err = os.Stat(fmt.Sprintf("%s/amnesia", userDir))
	if os.IsNotExist(err) {
		err = os.Mkdir(fmt.Sprintf("%s/amnesia", userDir), 0755)
		if err != nil {
			return s, err
		}
	}

	file, err := json.Marshal(s)
	if err != nil {
		return s, err
	}

	err = os.WriteFile(fmt.Sprintf("%s/amnesia/config.json", userDir), file, 0644)
	if err != nil {
		return s, err
	}

	return s, err
}

func LoadUserSettings() (s UserSettings, err error) {
	userDir := GetSaveDir()

	// TLDR: this has to stay config.json for compatibility reasons.
	_, err = os.Stat(fmt.Sprintf("%s/amnesia/config.json", userDir))
	if os.IsNotExist(err) {
		return NewUserSettings()
	}

	file, err := os.ReadFile(fmt.Sprintf("%s/amnesia/config.json", userDir))
	if err != nil {
		return s, err
	}

	err = json.Unmarshal([]byte(file), &s)
	if err != nil {
		return s, err
	}

	return s, err
}

func (s *UserSettings) SaveUserSettings() (err error) {

	userDir := GetSaveDir()
	file, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%s/amnesia/config.json", userDir), file, 0644)
	if err != nil {
		return err
	}

	return err
}
