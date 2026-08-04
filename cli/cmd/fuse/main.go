package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	fusecli "fuse/cli/internal"
	"fuse/cli/internal/auth"
	"fuse/cli/internal/ui"

	cli "github.com/urfave/cli/v2"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := fusecli.New(fusecli.DefaultOptions())
	if err == nil {
		err = app.RunContext(ctx, os.Args)
	}
	if err == nil {
		return
	}

	var presentedErr *ui.PresentedError
	if !errors.As(err, &presentedErr) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	var exitCoder cli.ExitCoder
	if errors.As(err, &exitCoder) {
		os.Exit(exitCoder.ExitCode())
	}
	if errors.Is(err, context.Canceled) {
		os.Exit(130)
	}
	if errors.Is(err, auth.ErrNotAuthenticated) {
		os.Exit(3)
	}
	os.Exit(1)
}
