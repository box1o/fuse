package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fuse/cli/internal/auth"
)

const (
	keyringService = "fuse"
	keyringAccount = "default"
)

type Store struct {
	fallbackPath string
	secretTool   string
}

func NewStore() (*Store, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve config directory: %w", err)
	}

	store := &Store{fallbackPath: filepath.Join(base, "fuse", "credentials.json")}
	if os.Getenv("DBUS_SESSION_BUS_ADDRESS") != "" {
		store.secretTool, _ = exec.LookPath("secret-tool")
	}

	return store, nil
}

func (s *Store) Save(ctx context.Context, credential auth.Credential) error {
	if credential.AccessToken == "" {
		return auth.ErrNotAuthenticated
	}

	payload, err := json.Marshal(credential)
	if err != nil {
		return fmt.Errorf("encode credential: %w", err)
	}
	if s.saveKeyring(ctx, payload) == nil && s.secretTool != "" {
		if err := os.Remove(s.fallbackPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove fallback credential: %w", err)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(s.fallbackPath), 0o700); err != nil {
		return fmt.Errorf("create credential directory: %w", err)
	}
	if err := os.Chmod(filepath.Dir(s.fallbackPath), 0o700); err != nil {
		return fmt.Errorf("protect credential directory: %w", err)
	}
	if err := writeProtectedFile(s.fallbackPath, payload); err != nil {
		return fmt.Errorf("write credential: %w", err)
	}
	return nil
}

func (s *Store) Load(ctx context.Context) (auth.Credential, error) {
	payload, err := s.loadKeyring(ctx)
	if err != nil {
		payload, err = os.ReadFile(s.fallbackPath)
	}
	if errors.Is(err, os.ErrNotExist) {
		return auth.Credential{}, auth.ErrNotAuthenticated
	}
	if err != nil {
		return auth.Credential{}, fmt.Errorf("read credential: %w", err)
	}

	var credential auth.Credential
	if err := json.Unmarshal(payload, &credential); err != nil {
		return auth.Credential{}, fmt.Errorf("decode credential: %w", err)
	}
	if credential.AccessToken == "" {
		return auth.Credential{}, auth.ErrNotAuthenticated
	}

	return credential, nil
}

func (s *Store) Delete(ctx context.Context) error {
	var keyringErr error
	if s.secretTool != "" {
		command := exec.CommandContext(ctx, s.secretTool, "clear", "service", keyringService, "account", keyringAccount)
		if err := command.Run(); err != nil && !isMissingKeyringEntry(err) {
			keyringErr = fmt.Errorf("clear keyring credential: %w", err)
		}
	}

	fileErr := os.Remove(s.fallbackPath)
	if errors.Is(fileErr, os.ErrNotExist) {
		fileErr = nil
	}
	if fileErr != nil {
		fileErr = fmt.Errorf("remove credential file: %w", fileErr)
	}

	return errors.Join(keyringErr, fileErr)
}

func writeProtectedFile(path string, payload []byte) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".credentials-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(payload); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

func isMissingKeyringEntry(err error) bool {
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && exitErr.ExitCode() == 1
}

func (s *Store) saveKeyring(ctx context.Context, payload []byte) error {
	if s.secretTool == "" {
		return errors.New("system keyring unavailable")
	}

	command := exec.CommandContext(ctx, s.secretTool, "store", "--label=Fuse CLI", "service", keyringService, "account", keyringAccount)
	command.Stdin = strings.NewReader(string(payload))
	return command.Run()
}

func (s *Store) loadKeyring(ctx context.Context) ([]byte, error) {
	if s.secretTool == "" {
		return nil, errors.New("system keyring unavailable")
	}

	command := exec.CommandContext(ctx, s.secretTool, "lookup", "service", keyringService, "account", keyringAccount)
	payload, err := command.Output()
	if err != nil || strings.TrimSpace(string(payload)) == "" {
		return nil, auth.ErrNotAuthenticated
	}

	return payload, nil
}
