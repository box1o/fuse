package interactive

import (
	"context"
	"errors"
	"io"

	"fuse/cli/internal/auth"
	"fuse/cli/internal/ui"

	"github.com/charmbracelet/huh"
)

type App struct {
	auth       *auth.Service
	prompt     *ui.Prompt
	presenter  ui.Presenter
	observer   auth.LoginObserver
	clientName string
}

func New(authService *auth.Service, input io.Reader, output io.Writer, mode ui.Mode, accessible bool) *App {
	return &App{
		auth:       authService,
		prompt:     ui.NewPrompt(input, output, accessible, mode.NoColor),
		presenter:  ui.NewPresenter(output, mode),
		observer:   ui.NewLoginObserver(output, mode),
		clientName: "Fuse CLI",
	}
}

func (a *App) Run(ctx context.Context) error {
	for {
		action, err := a.prompt.Select(ctx, "What do you want to do?", huh.NewOption("Authentication", "auth"), huh.NewOption("Exit", "exit"))
		if err != nil {
			return normalizeCancellation(err)
		}

		switch action {
		case "auth":
			if err := a.runAuth(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					return err
				}
				_ = a.presenter.Error(err)
			}
		case "exit":
			return nil
		}
	}
}

func (a *App) runAuth(ctx context.Context) error {
	action, err := a.prompt.Select(ctx, "Authentication", huh.NewOption("Login", "login"), huh.NewOption("Status", "status"), huh.NewOption("Logout", "logout"), huh.NewOption("Back", "back"))
	if err != nil {
		return normalizeCancellation(err)
	}

	switch action {
	case "login":
		session, err := a.auth.Login(ctx, auth.LoginOptions{ClientName: a.clientName, OpenBrowser: true, Observer: a.observer})
		if err != nil {
			return err
		}
		return a.presenter.Authenticated(session)
	case "status":
		status, err := a.auth.Status(ctx)
		if err != nil {
			return err
		}
		return a.presenter.AuthStatus(status)
	case "logout":
		confirmed, err := a.prompt.Confirm(ctx, "Log out from Fuse CLI?", false)
		if err != nil {
			return normalizeCancellation(err)
		}
		if !confirmed {
			return nil
		}
		if err := a.auth.Logout(ctx); err != nil {
			return err
		}
		return a.presenter.LoggedOut()
	case "back":
		return nil
	default:
		return nil
	}
}

func normalizeCancellation(err error) error {
	if errors.Is(err, huh.ErrUserAborted) {
		return nil
	}
	return err
}
