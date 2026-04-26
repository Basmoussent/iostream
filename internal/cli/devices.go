package cli

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/0xBasmoussent/iphone-mirror/internal/usb"
)

type devicesCmd struct{}

func (devicesCmd) Name() string    { return "devices" }
func (devicesCmd) Summary() string { return "List connected iOS devices." }

func (devicesCmd) Run(args []string, env Env) int {
	fs := newFlagSet("devices", env)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON instead of a table")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	devs, err := usb.NewBackend().Discover()
	if err != nil {
		fmt.Fprintf(env.Stderr, "iphone-mirror: %v\n", err)
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
		fmt.Fprintf(env.Stderr, "iphone-mirror: %v\n", err)
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
	fmt.Fprintln(tw, "UDID\tPRODUCT\tSERIAL\tBUS\tADDR\tQT")
	for _, d := range devs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%t\n",
			truncate(d.UDID, 12), d.Product, d.SerialNumber, d.BusNumber, d.DeviceAddress, d.QuickTimeEnabled)
	}
	if err := tw.Flush(); err != nil {
		fmt.Fprintf(env.Stderr, "iphone-mirror: %v\n", err)
		return 1
	}
	return 0
}

// truncate shortens a string to n runes, appending an ellipsis when it had to cut.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
}
