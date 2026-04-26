package cli

import (
	"fmt"

	"github.com/Basmoussent/iostream/internal/usb"
)

type activateCmd struct{}

func (activateCmd) Name() string { return "activate" }
func (activateCmd) Summary() string {
	return "Switch an iPhone into Apple's QuickTime USB configuration."
}

func (activateCmd) Run(args []string, env Env) int {
	fs := newFlagSet("activate", env)
	udid := fs.String("udid", "", "device UDID (USB serial). Empty = first device found.")
	if cont, code := parseFlags(fs, args); !cont {
		return code
	}

	if err := usb.NewBackend().Activate(*udid); err != nil {
		fmt.Fprintf(env.Stderr, "iostream: %v\n", err)
		return 1
	}
	fmt.Fprintln(env.Stdout, "QuickTime configuration activated.")
	return 0
}

type deactivateCmd struct{}

func (deactivateCmd) Name() string    { return "deactivate" }
func (deactivateCmd) Summary() string { return "Restore an iPhone's default USB configuration." }

func (deactivateCmd) Run(args []string, env Env) int {
	fs := newFlagSet("deactivate", env)
	udid := fs.String("udid", "", "device UDID (USB serial). Empty = first device found.")
	if cont, code := parseFlags(fs, args); !cont {
		return code
	}

	if err := usb.NewBackend().Deactivate(*udid); err != nil {
		fmt.Fprintf(env.Stderr, "iostream: %v\n", err)
		return 1
	}
	fmt.Fprintln(env.Stdout, "QuickTime configuration deactivated.")
	return 0
}
