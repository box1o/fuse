package fusecli

import (
	"io"
	"os"

	"fuse/cli/internal/platform"
)

type Options struct {
	Input       io.Reader
	Output      io.Writer
	ErrorOutput io.Writer
	Interactive func() bool
}

func DefaultOptions() Options {
	return Options{
		Input:       os.Stdin,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
		Interactive: platform.Interactive,
	}
}

func (o Options) withDefaults() Options {
	defaults := DefaultOptions()
	if o.Input == nil {
		o.Input = defaults.Input
	}
	if o.Output == nil {
		o.Output = defaults.Output
	}
	if o.ErrorOutput == nil {
		o.ErrorOutput = defaults.ErrorOutput
	}
	if o.Interactive == nil {
		o.Interactive = defaults.Interactive
	}
	return o
}
