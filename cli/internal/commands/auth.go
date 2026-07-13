package commands

import (
	"fuse/cli/internal/application"
	"fuse/cli/internal/ports"

	cli "github.com/urfave/cli/v2"
)

type authCommands struct {
	service *application.AuthService
	output  ports.Output
}

func authCommand(service *application.AuthService, output ports.Output) *cli.Command {
	commands := &authCommands{service: service, output: output}
	return &cli.Command{
		Name:  "auth",
		Usage: "manage CLI authentication",
		Subcommands: []*cli.Command{
			{
				Name:   "login",
				Usage:  "authenticate using a browser device code",
				Action: commands.login,
			},
			{
				Name:   "status",
				Usage:  "show the authenticated account",
				Action: commands.status,
			},
			{
				Name:   "logout",
				Usage:  "revoke and remove the CLI credential",
				Action: commands.logout,
			},
		},
	}
}

func (c *authCommands) login(ctx *cli.Context) error {
	_, err := c.service.Login(contextOf(ctx))
	return err
}

func (c *authCommands) status(ctx *cli.Context) error {
	status, err := c.service.Status(contextOf(ctx))
	if err != nil {
		return err
	}
	return c.output.JSON(status)
}

func (c *authCommands) logout(ctx *cli.Context) error {
	if err := c.service.Logout(contextOf(ctx)); err != nil {
		return err
	}
	c.output.Println("Logged out from Fuse CLI.")
	return nil
}
