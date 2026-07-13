package commands

import (
	"fmt"
	"strings"

	"fuse/cli/internal/application"
	"fuse/cli/internal/ports"
	"fuse/pkg/computeapi"

	cli "github.com/urfave/cli/v2"
)

type computeCommands struct {
	service *application.ComputeService
	output  ports.Output
	prompt  ports.Prompter
}

func computeCommand(service *application.ComputeService, output ports.Output, prompt ports.Prompter) *cli.Command {
	commands := &computeCommands{service: service, output: output, prompt: prompt}
	return &cli.Command{
		Name:  "compute",
		Usage: "inspect and manage compute nodes",
		Subcommands: []*cli.Command{
			commands.inspect(),
			commands.register(),
			commands.list(),
			commands.get(),
			commands.update(),
			commands.delete(),
		},
	}
}

func (c *computeCommands) inspect() *cli.Command {
	return &cli.Command{
		Name:   "inspect",
		Usage:  "inspect local compute capabilities",
		Flags:  []cli.Flag{&cli.BoolFlag{Name: "json"}},
		Action: c.runInspect,
	}
}

func (c *computeCommands) register() *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "register this machine",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name"},
			&cli.BoolFlag{Name: "yes", Aliases: []string{"y"}},
		},
		Action: c.runRegister,
	}
}

func (c *computeCommands) list() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "list registered nodes",
		Flags:  []cli.Flag{&cli.BoolFlag{Name: "json"}},
		Action: c.runList,
	}
}

func (c *computeCommands) get() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "get a node",
		ArgsUsage: "NODE_ID",
		Action:    c.runGet,
	}
}

func (c *computeCommands) update() *cli.Command {
	return &cli.Command{
		Name:      "update",
		Usage:     "rename, enable, or disable a node",
		ArgsUsage: "NODE_ID",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name"},
			&cli.BoolFlag{Name: "enable"},
			&cli.BoolFlag{Name: "disable"},
		},
		Action: c.runUpdate,
	}
}

func (c *computeCommands) delete() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a node",
		ArgsUsage: "NODE_ID",
		Flags:     []cli.Flag{&cli.BoolFlag{Name: "yes", Aliases: []string{"y"}}},
		Action:    c.runDelete,
	}
}

func (c *computeCommands) runInspect(ctx *cli.Context) error {
	capabilities, err := c.service.Inspect(contextOf(ctx))
	if err != nil {
		return err
	}
	if ctx.Bool("json") {
		return c.output.JSON(capabilities)
	}
	c.output.Capabilities(capabilities)
	return nil
}

func (c *computeCommands) runRegister(ctx *cli.Context) error {
	capabilities, err := c.service.Inspect(contextOf(ctx))
	if err != nil {
		return err
	}
	c.output.Capabilities(capabilities)
	if !ctx.Bool("yes") && !c.prompt.Confirm("Register this machine? [y/N] ") {
		return fmt.Errorf("registration cancelled")
	}

	node, created, err := c.service.Register(contextOf(ctx), strings.TrimSpace(ctx.String("name")), capabilities)
	if err != nil {
		return err
	}
	verb := "Updated"
	if created {
		verb = "Registered"
	}
	c.output.Printf("%s compute node %s (%s).\n", verb, node.Name, node.ID)
	return nil
}

func (c *computeCommands) runList(ctx *cli.Context) error {
	nodes, err := c.service.List(contextOf(ctx))
	if err != nil {
		return err
	}
	if ctx.Bool("json") {
		return c.output.JSON(nodes)
	}
	if len(nodes) == 0 {
		c.output.Println("No compute nodes registered.")
		return nil
	}
	c.output.Printf("%-36s  %-24s  %-12s  %s\n", "ID", "NAME", "STATUS", "HOSTNAME")
	for _, node := range nodes {
		c.output.Printf("%-36s  %-24s  %-12s  %s\n", node.ID, node.Name, node.Status, node.Hostname)
	}
	return nil
}

func (c *computeCommands) runGet(ctx *cli.Context) error {
	if err := requireOneArgument(ctx, "NODE_ID"); err != nil {
		return err
	}
	node, err := c.service.Get(contextOf(ctx), ctx.Args().First())
	if err != nil {
		return err
	}
	return c.output.JSON(node)
}

func (c *computeCommands) runUpdate(ctx *cli.Context) error {
	if err := requireOneArgument(ctx, "NODE_ID"); err != nil {
		return err
	}
	if ctx.Bool("enable") && ctx.Bool("disable") {
		return cli.Exit("--enable and --disable cannot be combined", 2)
	}

	request := computeapi.UpdateNodeRequest{}
	if name := strings.TrimSpace(ctx.String("name")); name != "" {
		request.Name = &name
	}
	if ctx.Bool("enable") {
		disabled := false
		request.Disabled = &disabled
	}
	if ctx.Bool("disable") {
		disabled := true
		request.Disabled = &disabled
	}

	node, err := c.service.Update(contextOf(ctx), ctx.Args().First(), request)
	if err != nil {
		return err
	}
	c.output.Printf("Updated compute node %s (%s).\n", node.Name, node.ID)
	return nil
}

func (c *computeCommands) runDelete(ctx *cli.Context) error {
	if err := requireOneArgument(ctx, "NODE_ID"); err != nil {
		return err
	}
	nodeID := ctx.Args().First()
	if !ctx.Bool("yes") && !c.prompt.Confirm("Delete compute node "+nodeID+"? [y/N] ") {
		return fmt.Errorf("deletion cancelled")
	}
	if err := c.service.Delete(contextOf(ctx), nodeID); err != nil {
		return err
	}
	c.output.Println("Compute node deleted.")
	return nil
}

func requireOneArgument(ctx *cli.Context, name string) error {
	if ctx.NArg() != 1 {
		return cli.Exit(name+" is required", 2)
	}
	return nil
}
