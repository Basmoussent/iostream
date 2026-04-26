package cli

import (
	"errors"
	"fmt"

	"github.com/Basmoussent/iostream/internal/driver"
)

type setupDriverCmd struct{}

func (setupDriverCmd) Name() string { return "setup-driver" }
func (setupDriverCmd) Summary() string {
	return "Install the WinUSB driver on the iPhone interface (Windows only)."
}

func (setupDriverCmd) Run(args []string, env Env) int {
	fs := newFlagSet("setup-driver", env)
	if cont, code := parseFlags(fs, args); !cont {
		return code
	}

	if err := driver.Setup(); err != nil {
		if errors.Is(err, driver.ErrUnsupportedPlatform) {
			fmt.Fprintln(env.Stdout, "iostream: no driver install required on this OS — skipping.")
			return 0
		}
		fmt.Fprintf(env.Stderr, "iostream: %v\n", err)
		return 1
	}
	fmt.Fprintln(env.Stdout, "iostream: driver setup complete. Run 'iostream devices' to confirm.")
	return 0
}
