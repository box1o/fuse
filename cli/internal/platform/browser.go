package platform

import (
	"fmt"
	"os/exec"
	"runtime"
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
		return fmt.Errorf("opening a browser is unsupported on %s", runtime.GOOS)
	}

	if err := command.Start(); err != nil {
		return fmt.Errorf("open browser: %w", err)
	}
	go func() {
		_ = command.Wait()
	}()
	return nil
}
