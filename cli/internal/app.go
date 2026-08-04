package fusecli

import (
	"fmt"

	"fuse/cli/internal/commands"
	"fuse/cli/internal/config"

	cli "github.com/urfave/cli/v2"
)

func New(options Options) (*cli.App, error) {
	options = options.withDefaults()
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load CLI configuration: %w", err)
	}
	dependencies, err := buildDependencies(options, cfg)
	if err != nil {
		return nil, err
	}

	return commands.New(dependencies, cfg), nil
}
