package bootstrap

import (
	"fmt"
	"io"

	"fuse/cli/internal/application"
	"fuse/cli/internal/commands"
	"fuse/cli/internal/infrastructure/capabilities"
	"fuse/cli/internal/infrastructure/httpapi"
	"fuse/cli/internal/infrastructure/storage"
	"fuse/cli/internal/infrastructure/system"
	"fuse/cli/internal/presentation"

	cli "github.com/urfave/cli/v2"
)

func Build(stdin io.Reader, stdout io.Writer) (*cli.App, error) {
	settings, err := storage.NewSettingsManager()
	if err != nil {
		return nil, fmt.Errorf("load CLI settings: %w", err)
	}
	credentials, err := storage.NewCredentialStore()
	if err != nil {
		return nil, fmt.Errorf("initialize credential storage: %w", err)
	}

	console := presentation.NewConsole(stdin, stdout)
	httpClient := httpapi.NewClient(settings.Current().APIURL)
	authService := application.NewAuthService(
		httpapi.NewDeviceAuthGateway(httpClient), credentials, system.Browser{}, system.Waiter{}, console,
	)
	computeService := application.NewComputeService(
		httpapi.NewNodeGateway(httpClient), credentials, settings, capabilities.NewLinuxDetector(),
	)
	return commands.New(
		authService,
		computeService,
		console,
		console,
		httpClient,
		settings.Current().APIURL,
	), nil
}
