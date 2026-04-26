// Package usb describes iOS devices reachable over USB and the contract that
// any platform-specific backend must satisfy.
//
// Two backends ship in-tree:
//
//   - The libusb backend (build tag `cgo`) talks to the device for real via
//     gousb + Apple's hidden QuickTime USB configuration.
//   - The stub backend (build tag `!cgo`) compiles everywhere and returns
//     ErrBackendUnavailable from every operation, which the CLI surfaces as a
//     friendly "no backend in this build" message.
//
// The stub exists so the project keeps building cleanly on hosts without
// libusb headers — useful for tooling, CI lint jobs, and people who only want
// to read the source.
package usb

import (
	"errors"
	"io"
)

// ErrBackendUnavailable is returned when no real USB backend is compiled in.
// Callers should detect it and surface a helpful message rather than treating
// it as a hard failure.
var ErrBackendUnavailable = errors.New("usb: no libusb backend compiled in (build with CGO_ENABLED=1, see docs/WINDOWS.md)")

// ErrDeviceNotFound is returned when a UDID does not match any attached device.
var ErrDeviceNotFound = errors.New("usb: device not found")

// Backend is the full set of operations the CLI needs from a USB driver.
//
// Implementations must be safe to use from a single goroutine; the CLI never
// fans out work onto multiple goroutines for the same backend instance.
type Backend interface {
	// Discover lists every Apple device currently visible on the USB bus.
	Discover() ([]Device, error)

	// Activate flips the device into Apple's hidden QuickTime USB configuration.
	// The device disappears from the bus for a few hundred ms and reappears
	// with the H.264 stream endpoint exposed; implementations must handle the
	// re-enumeration internally before returning.
	Activate(udid string) error

	// Deactivate puts the device back into its default USB configuration.
	Deactivate(udid string) error

	// Stream consumes the QuickTime stream from the device and writes:
	//   - H.264 NAL units in Annex-B framing to opts.Video
	//   - linear PCM audio samples to opts.Audio
	// It blocks until opts.Stop is closed or the device disconnects.
	Stream(udid string, opts StreamOptions) error
}

// StreamOptions configures a streaming session.
//
// Either Video or Audio may be nil to suppress that track. AutoActivate runs
// Activate(udid) before reading from the stream; the CLI uses it so that the
// happy path is `iphone-mirror stream | ffplay -` with nothing else to set up.
type StreamOptions struct {
	Video        io.Writer
	Audio        io.Writer
	Stop         <-chan struct{}
	AutoActivate bool
}

// NewBackend returns the backend wired in by build tags. See package doc.
//
// The actual factory lives in backend_libusb.go (cgo) or backend_stub.go
// (!cgo). Keeping the declaration here lets callers depend on a stable
// symbol regardless of how the binary was built.
func NewBackend() Backend {
	return newBackend()
}
