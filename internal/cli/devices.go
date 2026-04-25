package cli

import (
	"errors"
	"fmt"
)

type devicesCmd struct{}

func (devicesCmd) Name() string    { return "devices" }
func (devicesCmd) Summary() string { return "List connected iOS devices." }

func (devicesCmd) Run(args []string, env Env) int {
	fs := newFlagSet("devices", env)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON instead of a table")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, errPrintedHelp) {
			return 0
		}
		return 2
	}

	// USB discovery is wired in once internal/usb lands. For now, fail loudly
	// rather than silently print an empty list.
	_ = jsonOut
	fmt.Fprintln(env.Stderr, "iphone-mirror: device discovery not yet wired up — see docs/ROADMAP.md")
	return 1
}

// errPrintedHelp is returned from a FlagSet when the user passed -h or --help
// and the FlagSet already wrote the usage. We surface it as a sentinel so
// commands can distinguish "the user asked for help" from "the flags are bad".
var errPrintedHelp = errors.New("flag: help requested")
