package storage

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CredentialStore struct{ fallbackPath string }

func NewCredentialStore() (*CredentialStore, error) {
	directory, err := configDirectoryPath()
	if err != nil {
		return nil, err
	}
	return &CredentialStore{fallbackPath: filepath.Join(directory, "credentials")}, nil
}

func (s *CredentialStore) Save(token string) error {
	if path, ok := secretTool(); ok {
		command := exec.Command(path, "store", "--label=Fuse CLI", "service", "fuse", "account", "default")
		command.Stdin = strings.NewReader(token)
		if err := command.Run(); err == nil {
			_ = os.Remove(s.fallbackPath)
			return nil
		}
	}
	fmt.Fprintln(os.Stderr, "warning: system keyring unavailable; storing the CLI credential in a 0600 file")
	if err := ensureConfigDirectory(); err != nil {
		return err
	}
	return os.WriteFile(s.fallbackPath, []byte(token), 0o600)
}

func (s *CredentialStore) Load() (string, error) {
	if path, ok := secretTool(); ok {
		output, err := exec.Command(path, "lookup", "service", "fuse", "account", "default").Output()
		if err == nil && strings.TrimSpace(string(output)) != "" {
			return strings.TrimSpace(string(output)), nil
		}
	}
	payload, err := os.ReadFile(s.fallbackPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(payload)), nil
}

func (s *CredentialStore) Delete() error {
	if path, ok := secretTool(); ok {
		_ = exec.Command(path, "clear", "service", "fuse", "account", "default").Run()
	}
	err := os.Remove(s.fallbackPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func secretTool() (string, bool) {
	if os.Getenv("DBUS_SESSION_BUS_ADDRESS") == "" {
		return "", false
	}
	path, err := exec.LookPath("secret-tool")
	return path, err == nil
}
