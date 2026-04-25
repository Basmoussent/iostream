// Package usb describes iOS devices reachable over USB and the contract that
// any platform-specific discovery backend must satisfy.
//
// The actual libusb-backed implementation lives in a sibling file (added in a
// later milestone). Keeping the types and interface in their own file means
// the rest of the codebase can depend on stable shapes today.
package usb

// Device is one connected iOS device as seen from the host's USB stack.
//
// UDID is empty until the device has been queried via lockdownd (planned).
// QuickTimeEnabled reports whether the hidden QuickTime USB configuration is
// currently active — when true the device is exposing the H.264 video stream
// and we can start consuming it.
type Device struct {
	UDID             string `json:"udid"`
	Product          string `json:"product"`
	SerialNumber     string `json:"serial_number"`
	BusNumber        int    `json:"bus_number"`
	DeviceAddress    int    `json:"device_address"`
	QuickTimeEnabled bool   `json:"quicktime_enabled"`
}

// Discoverer enumerates iOS devices currently attached over USB.
type Discoverer interface {
	Discover() ([]Device, error)
}
