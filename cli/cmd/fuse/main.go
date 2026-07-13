package main

import (
	"context"
	"fmt"
	"os"

	"fuse/cli/internal/bootstrap"
)

func main() {
	app, err := bootstrap.Build(os.Stdin, os.Stdout)
	if err == nil {
		err = app.RunContext(context.Background(), os.Args)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
