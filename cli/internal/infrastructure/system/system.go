package system

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

type Browser struct{}

func (Browser) Open(url string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		command = exec.Command("xdg-open", url)
	case "darwin":
		command = exec.Command("open", url)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform %s", runtime.GOOS)
	}
	return command.Start()
}

type Waiter struct{}

func (Waiter) Wait(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
