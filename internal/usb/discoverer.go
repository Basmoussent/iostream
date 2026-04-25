package usb

import "errors"

// ErrBackendUnavailable is returned by the default discoverer until a real
// libusb-backed implementation is wired in. Callers should surface it as a
// "not yet implemented" message rather than treating it as a hard failure.
var ErrBackendUnavailable = errors.New("usb: discovery backend not yet implemented (see docs/ROADMAP.md)")

// NewDiscoverer returns the discoverer for the current platform. Today every
// platform falls back to the stub; the libusb backend will register itself
// here once it lands.
func NewDiscoverer() Discoverer {
	return stubDiscoverer{}
}

type stubDiscoverer struct{}

func (stubDiscoverer) Discover() ([]Device, error) {
	return nil, ErrBackendUnavailable
}
