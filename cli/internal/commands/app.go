package commands

import (
	"context"

	"fuse/cli/internal/application"
	"fuse/cli/internal/ports"

	cli "github.com/urfave/cli/v2"
)

func New(
	auth *application.AuthService,
	compute *application.ComputeService,
	output ports.Output,
	prompt ports.Prompter,
	endpoint ports.APIEndpoint,
	defaultAPIURL string,
) *cli.App {
	app := cli.NewApp()
	app.Name = "fuse"
	app.Usage = "authenticate and manage Fuse compute nodes"
	app.Version = application.AgentVersion
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "api-url",
			Usage:   "Fuse server API URL",
			Value:   defaultAPIURL,
			EnvVars: []string{"FUSE_API_URL"},
		},
	}
	app.Before = func(c *cli.Context) error {
		return endpoint.SetBaseURL(c.String("api-url"))
	}
	app.Commands = []*cli.Command{
		authCommand(auth, output),
		computeCommand(compute, output, prompt),
	}
	return app
}

func contextOf(c *cli.Context) context.Context { return c.Context }
