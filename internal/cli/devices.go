package cli

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/Basmoussent/iostream/internal/usb"
)

type devicesCmd struct{}

func (devicesCmd) Name() string    { return "devices" }
func (devicesCmd) Summary() string { return "List connected iOS devices." }

func (devicesCmd) Run(args []string, env Env) int {
	fs := newFlagSet("devices", env)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON instead of a table")
	if cont, code := parseFlags(fs, args); !cont {
		return code
	}

	devs, err := usb.NewBackend().Discover()
	if err != nil {
		fmt.Fprintf(env.Stderr, "iostream: %v\n", err)
		return 1
	}

	if *jsonOut {
		return renderDevicesJSON(env, devs)
	}
	return renderDevicesTable(env, devs)
}

func renderDevicesJSON(env Env, devs []usb.Device) int {
	enc := json.NewEncoder(env.Stdout)
	enc.SetIndent("", "  ")
	if devs == nil {
		devs = []usb.Device{}
	}
	if err := enc.Encode(devs); err != nil {
		fmt.Fprintf(env.Stderr, "iostream: %v\n", err)
		return 1
	}
	return 0
}

func renderDevicesTable(env Env, devs []usb.Device) int {
	if len(devs) == 0 {
		fmt.Fprintln(env.Stdout, "No iOS devices found.")
		return 0
	}
	tw := tabwriter.NewWriter(env.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "UDID\tPRODUCT\tUSB\tQT")
	for _, d := range devs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%t\n",
			d.UDID, d.ProductName, d.USBInfo, d.QuickTimeEnabled)
	}
	if err := tw.Flush(); err != nil {
		fmt.Fprintf(env.Stderr, "iostream: %v\n", err)
		return 1
	}
	return 0
}
