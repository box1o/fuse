package commands

import (
	"fuse/cli/internal/auth"
	"fuse/cli/internal/ui"

	cli "github.com/urfave/cli/v2"
)

func authCommand(dependencies Dependencies) *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "manage CLI authentication",
		Subcommands: []*cli.Command{
			{
				Name:  "login",
				Usage: "authenticate through the browser device flow",
				Flags: append(outputFlags(), &cli.BoolFlag{Name: "no-browser", Usage: "do not open the browser automatically"}),
				Action: func(ctx *cli.Context) error {
					return runAction(ctx, dependencies, func() error {
						if err := validateOutputFlags(ctx); err != nil {
							return err
						}
						mode := outputMode(ctx)
						observer := ui.NewLoginObserver(dependencies.Output, mode)
						session, err := dependencies.Auth.Login(contextOf(ctx), auth.LoginOptions{ClientName: "Fuse CLI", OpenBrowser: !ctx.Bool("no-browser"), Observer: observer})
						if err != nil {
							return err
						}
						return ui.NewPresenter(dependencies.Output, mode).Authenticated(session)
					})
				},
			},
			{
				Name:  "status",
				Usage: "show the authenticated account",
				Flags: outputFlags(),
				Action: func(ctx *cli.Context) error {
					return runAction(ctx, dependencies, func() error {
						if err := validateOutputFlags(ctx); err != nil {
							return err
						}
						status, err := dependencies.Auth.Status(contextOf(ctx))
						if err != nil {
							return err
						}
						return ui.NewPresenter(dependencies.Output, outputMode(ctx)).AuthStatus(status)
					})
				},
			},
			{
				Name:  "logout",
				Usage: "revoke and remove the CLI credential",
				Flags: append(outputFlags(), &cli.BoolFlag{Name: "yes", Aliases: []string{"y"}, Usage: "skip confirmation"}),
				Action: func(ctx *cli.Context) error {
					return runAction(ctx, dependencies, func() error {
						if err := validateOutputFlags(ctx); err != nil {
							return err
						}
						confirmed, err := confirmLogout(ctx, dependencies)
						if err != nil || !confirmed {
							return err
						}
						if err := dependencies.Auth.Logout(contextOf(ctx)); err != nil {
							return err
						}
						return ui.NewPresenter(dependencies.Output, outputMode(ctx)).LoggedOut()
					})
				},
			},
		},
	}
}
