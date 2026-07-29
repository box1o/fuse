package fusecli

import (
	"fmt"

	"fuse/cli/internal/api"
	"fuse/cli/internal/auth"
	"fuse/cli/internal/commands"
	"fuse/cli/internal/config"
	"fuse/cli/internal/credentials"
	"fuse/cli/internal/platform"
)

func buildDependencies(options Options, cfg config.Config) (commands.Dependencies, error) {
	client, err := api.NewClient(cfg.APIURL)
	if err != nil {
		return commands.Dependencies{}, fmt.Errorf("create API client: %w", err)
	}
	credentialStore, err := credentials.NewStore()
	if err != nil {
		return commands.Dependencies{}, fmt.Errorf("create credential store: %w", err)
	}

	authService := auth.NewService(api.NewAuthGateway(client), credentialStore, platform.Browser{}, platform.Waiter{})
	return commands.Dependencies{
		Auth:        authService,
		API:         client,
		Input:       options.Input,
		Output:      options.Output,
		ErrorOutput: options.ErrorOutput,
		Interactive: options.Interactive,
	}, nil
}
