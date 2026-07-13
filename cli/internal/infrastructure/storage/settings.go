package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"fuse/cli/internal/ports"

	"github.com/google/uuid"
)

type SettingsManager struct {
	mu       sync.RWMutex
	path     string
	settings ports.Settings
}

type settingsFile struct {
	APIURL         string `json:"api_url"`
	AppURL         string `json:"app_url"`
	InstallationID string `json:"installation_id,omitempty"`
}

func NewSettingsManager() (*SettingsManager, error) {
	directory, err := configDirectoryPath()
	if err != nil {
		return nil, err
	}
	manager := &SettingsManager{path: filepath.Join(directory, "config.json")}
	manager.settings = ports.Settings{
		APIURL: envOrDefault("FUSE_API_URL", "https://mback.teckstate.com"),
		AppURL: envOrDefault("FUSE_APP_URL", "https://app.teckstate.com"),
	}
	payload, err := os.ReadFile(manager.path)
	if err == nil {
		var file settingsFile
		if err := json.Unmarshal(payload, &file); err != nil {
			return nil, err
		}
		if os.Getenv("FUSE_API_URL") == "" {
			manager.settings.APIURL = file.APIURL
		}
		if os.Getenv("FUSE_APP_URL") == "" {
			manager.settings.AppURL = file.AppURL
		}
		manager.settings.InstallationID = file.InstallationID
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	manager.settings.APIURL = strings.TrimRight(manager.settings.APIURL, "/")
	manager.settings.AppURL = strings.TrimRight(manager.settings.AppURL, "/")
	return manager, nil
}

func (m *SettingsManager) Current() ports.Settings {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.settings
}

func (m *SettingsManager) EnsureInstallationID() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, err := uuid.Parse(m.settings.InstallationID); err == nil {
		return m.settings.InstallationID, nil
	}
	m.settings.InstallationID = uuid.NewString()
	file := settingsFile{APIURL: m.settings.APIURL, AppURL: m.settings.AppURL, InstallationID: m.settings.InstallationID}
	payload, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return "", err
	}
	if err := ensureConfigDirectory(); err != nil {
		return "", err
	}
	if err := os.WriteFile(m.path, payload, 0o600); err != nil {
		return "", err
	}
	return m.settings.InstallationID, nil
}

func configDirectoryPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "fuse"), nil
}

func ensureConfigDirectory() error {
	directory, err := configDirectoryPath()
	if err != nil {
		return err
	}
	return os.MkdirAll(directory, 0o700)
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
