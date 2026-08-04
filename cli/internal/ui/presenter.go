package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"fuse/cli/internal/auth"

	"github.com/charmbracelet/lipgloss"
)

type Mode struct {
	JSON    bool
	Quiet   bool
	NoColor bool
}

type Presenter interface {
	Authenticated(auth.Session) error
	AuthStatus(auth.Status) error
	LoggedOut() error
	Error(error) error
}

type PresentedError struct {
	err error
}

func (e *PresentedError) Error() string {
	return e.err.Error()
}

func (e *PresentedError) Unwrap() error {
	return e.err
}

func PresentError(presenter Presenter, value error) error {
	if value == nil {
		return nil
	}
	if err := presenter.Error(value); err != nil {
		return fmt.Errorf("present error: %w", err)
	}
	return &PresentedError{err: value}
}

func NewPresenter(output io.Writer, mode Mode) Presenter {
	if mode.Quiet {
		return quietPresenter{}
	}
	if mode.JSON {
		return &jsonPresenter{output: output}
	}
	return &humanPresenter{output: output, noColor: mode.NoColor}
}

type humanPresenter struct {
	output  io.Writer
	noColor bool
}

func (p *humanPresenter) Authenticated(session auth.Session) error {
	_, err := fmt.Fprintf(p.output, "\n%s  Authentication successful\n│  User  %s\n└  Ready\n", p.success("◇"), userLabel(session.User))
	return err
}

func (p *humanPresenter) AuthStatus(status auth.Status) error {
	_, err := fmt.Fprintf(p.output, "%s  Authenticated\n│  User     %s\n│  Expires  %s\n└  Connected\n", p.success("◇"), userLabel(status.User), formatTime(status.ExpiresAt))
	return err
}

func (p *humanPresenter) LoggedOut() error {
	_, err := fmt.Fprintf(p.output, "%s  Logged out from Fuse CLI\n", p.success("◇"))
	return err
}

func (p *humanPresenter) Error(value error) error {
	_, err := fmt.Fprintf(p.output, "%s  %s\n", p.failure("✕"), value)
	return err
}

func (p *humanPresenter) success(value string) string {
	if p.noColor {
		return value
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render(value)
}

func (p *humanPresenter) failure(value string) string {
	if p.noColor {
		return value
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(value)
}

type jsonPresenter struct {
	output io.Writer
}

func (p *jsonPresenter) Authenticated(session auth.Session) error {
	return p.write(map[string]any{"authenticated": true, "user": session.User, "expires_at": session.Credential.ExpiresAt})
}

func (p *jsonPresenter) AuthStatus(status auth.Status) error {
	return p.write(status)
}

func (p *jsonPresenter) LoggedOut() error {
	return p.write(map[string]any{"authenticated": false})
}

func (p *jsonPresenter) Error(value error) error {
	return p.write(map[string]any{"error": map[string]string{"message": value.Error()}})
}

func (p *jsonPresenter) write(value any) error {
	encoder := json.NewEncoder(p.output)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

type quietPresenter struct{}

func (quietPresenter) Authenticated(auth.Session) error { return nil }
func (quietPresenter) AuthStatus(auth.Status) error     { return nil }
func (quietPresenter) LoggedOut() error                 { return nil }
func (quietPresenter) Error(error) error                { return nil }

func userLabel(user auth.User) string {
	if user.Email != "" {
		return user.Email
	}
	if user.Name != "" {
		return user.Name
	}
	return user.ID
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return "unknown"
	}
	return value.Local().Format(time.RFC1123)
}
