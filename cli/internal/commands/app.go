package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"fuse/cli/internal/api"
	"fuse/cli/internal/auth"
	"fuse/cli/internal/config"
	"fuse/cli/internal/interactive"
	"fuse/cli/internal/ui"

	"github.com/charmbracelet/huh"
	cli "github.com/urfave/cli/v2"
)

var Version = "dev"

type Dependencies struct {
	Auth        *auth.Service
	API         *api.Client
	Input       io.Reader
	Output      io.Writer
	ErrorOutput io.Writer
	Interactive func() bool
}

func New(dependencies Dependencies, cfg config.Config) *cli.App {
	app := cli.NewApp()
	app.Name = "fuse"
	app.Usage = "authenticate and work with Fuse"
	app.Version = Version
	app.Writer = dependencies.Output
	app.ErrWriter = dependencies.ErrorOutput
	app.HideHelpCommand = true
	app.Flags = globalFlags(cfg)
	app.Before = func(ctx *cli.Context) error {
		if ctx.Bool("json") && ctx.Bool("quiet") {
			return cli.Exit("--json and --quiet cannot be combined", 2)
		}
		return dependencies.API.SetBaseURL(ctx.String("api-url"))
	}
	app.Action = func(ctx *cli.Context) error {
		if ctx.Bool("no-interactive") || dependencies.Interactive == nil || !dependencies.Interactive() {
			return cli.ShowAppHelp(ctx)
		}

		menu := interactive.New(dependencies.Auth, dependencies.Input, dependencies.Output, outputMode(ctx), os.Getenv("ACCESSIBLE") != "")
		return menu.Run(ctx.Context)
	}
	app.Commands = []*cli.Command{authCommand(dependencies)}
	return app
}

func globalFlags(cfg config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "api-url", Usage: "Fuse API URL", Value: cfg.APIURL, EnvVars: []string{config.EnvAPIURL}},
		&cli.BoolFlag{Name: "json", Usage: "print machine-readable JSON"},
		&cli.BoolFlag{Name: "quiet", Aliases: []string{"q"}, Usage: "suppress normal output"},
		&cli.BoolFlag{Name: "no-color", Usage: "disable colored output"},
		&cli.BoolFlag{Name: "no-interactive", Usage: "never prompt for input"},
	}
}

func outputFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{Name: "json", Usage: "print machine-readable JSON"},
		&cli.BoolFlag{Name: "quiet", Aliases: []string{"q"}, Usage: "suppress normal output"},
		&cli.BoolFlag{Name: "no-color", Usage: "disable colored output"},
		&cli.BoolFlag{Name: "no-interactive", Usage: "never prompt for input"},
	}
}

func outputMode(ctx *cli.Context) ui.Mode {
	return ui.Mode{JSON: boolFlag(ctx, "json"), Quiet: boolFlag(ctx, "quiet"), NoColor: boolFlag(ctx, "no-color")}
}

func validateOutputFlags(ctx *cli.Context) error {
	if boolFlag(ctx, "json") && boolFlag(ctx, "quiet") {
		return cli.Exit("--json and --quiet cannot be combined", 2)
	}
	return nil
}

func boolFlag(ctx *cli.Context, name string) bool {
	for _, current := range ctx.Lineage() {
		if current.IsSet(name) {
			return current.Bool(name)
		}
	}
	return ctx.Bool(name)
}

func runAction(ctx *cli.Context, dependencies Dependencies, action func() error) error {
	err := action()
	if err == nil {
		return nil
	}

	mode := outputMode(ctx)
	errorMode := mode
	if errorMode.Quiet {
		errorMode.Quiet = false
		errorMode.JSON = false
	}
	output := dependencies.ErrorOutput
	if errorMode.JSON {
		output = dependencies.Output
	}
	return ui.PresentError(ui.NewPresenter(output, errorMode), err)
}

func contextOf(ctx *cli.Context) context.Context {
	if ctx.Context == nil {
		return context.Background()
	}
	return ctx.Context
}

func confirmLogout(ctx *cli.Context, dependencies Dependencies) (bool, error) {
	if ctx.Bool("yes") {
		return true, nil
	}
	if boolFlag(ctx, "no-interactive") || dependencies.Interactive == nil || !dependencies.Interactive() {
		return false, cli.Exit("--yes is required in non-interactive mode", 2)
	}

	prompt := ui.NewPrompt(dependencies.Input, dependencies.Output, os.Getenv("ACCESSIBLE") != "", boolFlag(ctx, "no-color"))
	confirmed, err := prompt.Confirm(contextOf(ctx), "Log out from Fuse CLI?", false)
	if err != nil && err != huh.ErrUserAborted {
		return false, fmt.Errorf("confirm logout: %w", err)
	}
	return confirmed, nil
}
