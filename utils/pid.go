package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func NewPidFile() error {
	pid := os.Getpid()
	tempDir := os.TempDir()

	file, err := os.Create(fmt.Sprintf("%s/passbolt-gui/passbolt-gui.pid", tempDir))
	if err != nil {
		return err
	}
	file.WriteString(strconv.Itoa(pid))
	file.Close()

	return nil
}

func RemovePidFile() error {
	tempDir := os.TempDir()

	err := os.Remove(fmt.Sprintf("%s/passbolt-gui/passbolt-gui.pid", tempDir))
	if err != nil {
		return err
	}

	return nil
}

func CheckPidFile() error {
	tempDir := os.TempDir()

	_, err := os.Stat(fmt.Sprintf("%s/passbolt-gui/passbolt-gui.pid", tempDir))
	if !os.IsNotExist(err) {
		return errors.New("Passbolt GUI is already running")
	}

	return nil
}
