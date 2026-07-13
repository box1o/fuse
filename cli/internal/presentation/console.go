package presentation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"fuse/pkg/computeapi"
)

type Console struct {
	out io.Writer
	in  io.Reader
}

func NewConsole(in io.Reader, out io.Writer) *Console { return &Console{in: in, out: out} }

func (c *Console) Printf(format string, args ...any) { _, _ = fmt.Fprintf(c.out, format, args...) }
func (c *Console) Println(values ...any)             { _, _ = fmt.Fprintln(c.out, values...) }

func (c *Console) JSON(value any) error {
	encoder := json.NewEncoder(c.out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func (c *Console) Capabilities(value computeapi.Capabilities) {
	c.Printf("OS: %s %s (%s)\n", value.OS.Name, value.OS.Version, value.OS.Architecture)
	c.Printf("CPU: %s, %d physical / %d logical cores\n", value.CPU.Model, value.CPU.PhysicalCores, value.CPU.LogicalCores)
	c.Printf("Memory: %.2f GiB\n", float64(value.Memory.TotalBytes)/(1<<30))
	c.Printf("Storage: %.2f GiB\n", float64(value.Storage.TotalBytes)/(1<<30))
	c.Printf("Docker: %t %s\n", value.ContainerRuntime.Available, value.ContainerRuntime.Version)
	for _, accelerator := range value.Accelerators {
		c.Printf("Accelerator: %s %s (%.2f GiB)\n", accelerator.Vendor, accelerator.Model, float64(accelerator.MemoryBytes)/(1<<30))
	}
}

func (c *Console) Confirm(prompt string) bool {
	c.Printf("%s", prompt)
	value, _ := bufio.NewReader(c.in).ReadString('\n')
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "y" || value == "yes"
}

func (c *Console) ShowDeviceCode(userCode, verificationURI string) {
	c.Printf("First copy your one-time code: %s\n", userCode)
	c.Printf("Open %s in your browser to continue.\n", verificationURI)
}

func (c *Console) ShowAuthenticated(ownerName, ownerEmail string, expiresAt time.Time) {
	c.Printf("Logged in as %s <%s>. Credential expires %s.\n", ownerName, ownerEmail, expiresAt.Format(time.RFC3339))
}
