package ui

import (
	"encoding/json"
	"fmt"
	"io"

	"fuse/cli/internal/auth"

	"github.com/charmbracelet/lipgloss"
)

type LoginObserver struct {
	output io.Writer
	mode   Mode
}

func NewLoginObserver(output io.Writer, mode Mode) *LoginObserver {
	return &LoginObserver{output: output, mode: mode}
}

func (o *LoginObserver) DeviceCodeReady(code auth.DeviceCode) {
	if o.mode.Quiet {
		return
	}
	if o.mode.JSON {
		o.writeJSON(map[string]any{"event": "device_code", "user_code": code.UserCode, "verification_uri": code.VerificationURI, "expires_in": int(code.ExpiresIn.Seconds())})
		return
	}

	marker := "◇"
	if !o.mode.NoColor {
		marker = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render(marker)
	}
	fmt.Fprintf(o.output, "%s  Open this page:\n│  %s\n│\n│  Your one-time code:\n│  %s\n", marker, code.VerificationURI, code.UserCode)
}

func (o *LoginObserver) WaitingForAuthorization() {
	if o.mode.Quiet {
		return
	}
	if o.mode.JSON {
		o.writeJSON(map[string]string{"event": "waiting_for_authorization"})
		return
	}
	fmt.Fprintln(o.output, "│\n◇  Waiting for authorization...")
}

func (o *LoginObserver) writeJSON(value any) {
	_ = json.NewEncoder(o.output).Encode(value)
}
